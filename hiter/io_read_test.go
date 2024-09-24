package hiter_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"testing/iotest"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestDecode(t *testing.T) {
	type sample struct {
		V string
	}

	src := []byte(`
	{"V":"foo"}
	{"V":"bar"}
	{"V":"baz"}
	`)

	t.Run("ok", func(t *testing.T) {
		dec := json.NewDecoder(iotest.DataErrReader(bytes.NewReader(src)))

		result := hiter.Collect2(hiter.Decode[sample](dec))
		assert.DeepEqual(
			t,
			[]hiter.KeyValue[sample, error]{
				{K: sample{"foo"}},
				{K: sample{"bar"}},
				{K: sample{"baz"}},
			},
			result,
		)
	})
	t.Run("error", func(t *testing.T) {
		sampleErr := errors.New("sample")
		dec := json.NewDecoder(io.MultiReader(bytes.NewReader(src), iotest.ErrReader(sampleErr)))

		result := hiter.Collect2(
			hiter.LimitAfter2(
				func(_ sample, err error) bool { return err == nil },
				hiter.Decode[sample](dec),
			),
		)
		assert.DeepEqual(
			t,
			[]hiter.KeyValue[sample, error]{
				{K: sample{"foo"}},
				{K: sample{"bar"}},
				{K: sample{"baz"}},
				{V: sampleErr},
			},
			result,
			compareError,
		)
	})
}
