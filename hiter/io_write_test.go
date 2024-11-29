package hiter_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"slices"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestWrite(t *testing.T) {
	type testData struct {
		V string
	}
	src := []testData{{"foo"}, {"bar"}, {"baz"}}
	marshal := func(v testData, written int) ([]byte, error) {
		bin, err := json.Marshal(v)
		if err != nil {
			return bin, err
		}
		bin = append(bin, []byte("\n")...)
		// To be realistic,
		// prepend byte offset of start of next message,
		// and append byte offset of start of this message.
		return append(
			append(
				binary.BigEndian.AppendUint64(make([]byte, 0, 8), uint64(written+len(bin)+8)),
				bin...,
			),
			binary.BigEndian.AppendUint64(make([]byte, 0, 8), uint64(written))...,
		), nil
	}
	expected := [][]byte{
		// 0
		binary.BigEndian.AppendUint64(nil, 20),
		[]byte(`{"V":"foo"}` + "\n"), // 12
		binary.BigEndian.AppendUint64(nil, 0),
		// 12+8+8=28
		binary.BigEndian.AppendUint64(nil, 28+12+8),
		[]byte(`{"V":"bar"}` + "\n"), // 12
		binary.BigEndian.AppendUint64(nil, 28),
		// 28+12+8+8=56
		binary.BigEndian.AppendUint64(nil, 56+12+8),
		[]byte(`{"V":"baz"}` + "\n"),
		binary.BigEndian.AppendUint64(nil, 56),
	}

	t.Run("ok", func(t *testing.T) {
		var buf bytes.Buffer

		var (
			n   int
			err error
		)

		assertResult := func() {
			t.Helper()
			expected := hiter.AppendBytes(nil, slices.Values(expected))
			assert.Assert(
				t,
				cmp.DeepEqual(expected, buf.Bytes()),
				"diff = %s", goCmp.Diff(string(expected), buf.String()),
			)
			assert.NilError(t, err)
			assert.Equal(t, len(expected), n)
		}

		n, err = hiter.Write(
			&buf,
			marshal,
			slices.Values(src),
		)
		assertResult()

		buf.Reset()
		var keys []int
		n, err = hiter.Write2(
			&buf,
			func(k int, v testData, written int) ([]byte, error) {
				keys = append(keys, k)
				return marshal(v, written)
			},
			slices.All(src),
		)
		assertResult()
		assert.DeepEqual(t, slices.Collect(hiter.Range(0, 3)), keys)

	})

	t.Run("marshal error", func(t *testing.T) {
		var buf bytes.Buffer

		var (
			count int
			n     int
			err   error
		)

		assertResult := func() {
			t.Helper()
			expected := hiter.AppendBytes(nil, slices.Values(expected[:3]))
			assert.Assert(
				t,
				cmp.DeepEqual(expected, buf.Bytes()),
				"diff = %s", goCmp.Diff(string(expected), buf.String()),
			)
			assert.ErrorIs(t, err, testcase.ErrSample)
			assert.Equal(t, len(expected), n)
		}

		count = 0
		n, err = hiter.Write(
			&buf,
			func(v testData, written int) ([]byte, error) {
				if count == 1 {
					return []byte("wah"), testcase.ErrSample
				}
				count++
				return marshal(v, written)
			},
			slices.Values(src),
		)
		assertResult()

		count = 0
		buf.Reset()
		var keys []int
		n, err = hiter.Write2(
			&buf,
			func(k int, v testData, written int) ([]byte, error) {
				keys = append(keys, k)
				if count == 1 {
					return []byte("wah"), testcase.ErrSample
				}
				count++
				return marshal(v, written)
			},
			slices.All(src),
		)
		assertResult()
		assert.DeepEqual(t, slices.Collect(hiter.Range(0, 2)), keys)
	})

	t.Run("writer error", func(t *testing.T) {
		buf := new(bytes.Buffer)
		w := &errWriter{
			w:   buf,
			err: testcase.ErrSample,
		}

		var (
			n   int
			err error
		)
		assertResult := func() {
			t.Helper()
			expected := hiter.AppendBytes(nil, slices.Values(expected[:6]))
			assert.Assert(
				t,
				cmp.DeepEqual(expected, buf.Bytes()),
				"diff = %s", goCmp.Diff(string(expected), buf.String()),
			)
			assert.ErrorIs(t, err, testcase.ErrSample)
			assert.Equal(t, len(expected), n)
		}

		buf.Reset()
		w.n = 1
		n, err = hiter.Write(
			w,
			marshal,
			slices.Values(src),
		)
		assertResult()

		buf.Reset()
		w.n = 1
		var keys []int
		n, err = hiter.Write2(
			w,
			func(k int, v testData, written int) ([]byte, error) {
				keys = append(keys, k)
				return marshal(v, written)
			},
			slices.All(src),
		)
		assertResult()
		assert.DeepEqual(t, slices.Collect(hiter.Range(0, 3)), keys)
	})
}

type errWriter struct {
	w   io.Writer
	err error
	n   int
}

func (w *errWriter) Write(p []byte) (int, error) {
	if w.n < 0 {
		return 0, w.err
	}
	w.n--
	return w.w.Write(p)
}

func TestEncode(t *testing.T) {
	type testData struct {
		V string
	}
	src := []testData{{"foo"}, {"bar"}, {"baz"}}
	t.Run("ok", func(t *testing.T) {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		err := hiter.Encode(enc, slices.Values(src))
		assert.NilError(t, err)
		assert.DeepEqual(
			t,
			`{"V":"foo"}
{"V":"bar"}
{"V":"baz"}
`,
			buf.String(),
		)
	})
	t.Run("encoder error", func(t *testing.T) {
		buf := new(bytes.Buffer)
		w := &errWriter{
			w:   buf,
			err: testcase.ErrSample,
			n:   1,
		}
		enc := json.NewEncoder(w)
		err := hiter.Encode(enc, slices.Values(src))
		assert.ErrorIs(t, err, testcase.ErrSample)
		assert.DeepEqual(
			t,
			`{"V":"foo"}
{"V":"bar"}
`,
			buf.String(),
		)
	})
}
