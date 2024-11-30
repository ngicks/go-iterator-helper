package iterreader_test

import (
	"encoding/hex"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/ngicks/go-iterator-helper/hiter/cryptoiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/hiter/iterreader"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
)

func TestReader(t *testing.T) {
	src := slices.Collect(
		xiter.Limit(
			xiter.Map(
				func(b []byte) string { return hex.EncodeToString(b) },
				cryptoiter.RandBytes(8),
			),
			10,
		),
	)

	var (
		r   io.Reader
		err error
		bin []byte
	)

	r = iterreader.Reader(
		func(s string) ([]byte, error) { return []byte(s), nil },
		iterable.NewResumable(slices.Values(src)).IntoIter(),
	)
	bin, err = io.ReadAll(r)
	assert.NilError(t, err)
	assert.Equal(t, strings.Join(src, ""), string(bin))

	r = iterreader.Reader(
		func(s string) ([]byte, error) { return []byte(s + "\n"), nil },
		iterable.NewResumable(slices.Values(src)).IntoIter(),
	)
	bin, err = io.ReadAll(r)
	assert.NilError(t, err)
	assert.Equal(t, strings.Join(src, "\n")+"\n", string(bin))

	r = iterreader.Reader(
		func(s string) ([]byte, error) { return []byte(s), nil },
		iterable.NewResumable(slices.Values(src)).IntoIter(),
	)
	err = iotest.TestReader(r, []byte(strings.Join(src, "")))
	assert.NilError(t, err)
}

func TestMarshalerReader(t *testing.T) {
	times := []time.Time{
		time.Date(2024, 10, 20, 11, 20, 47, 0, time.UTC),
		time.Date(2024, 9, 23, 2, 45, 21, 0, time.UTC),
		time.Date(2024, 8, 4, 20, 1, 36, 0, time.UTC),
	}
	expectedStr := `2024-10-20T11:20:47Z
2024-09-23T02:45:21Z
2024-08-04T20:01:36Z
`
	expectedBytes := append(
		append(
			append(
				[]byte(nil),
				must(times[0].MarshalBinary())...,
			),
			must(times[1].MarshalBinary())...,
		),
		must(times[2].MarshalBinary())...,
	)

	var (
		bin []byte
		err error
	)
	bin, err = io.ReadAll(iterreader.TextMarshaler([]byte("\n"), slices.Values(times)))
	assert.NilError(t, err)
	assert.DeepEqual(t, expectedStr, string(bin))
	bin, err = io.ReadAll(iterreader.BinaryMarshaler(nil, slices.Values(times)))
	assert.NilError(t, err)
	assert.DeepEqual(t, expectedBytes, bin)
}

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}
