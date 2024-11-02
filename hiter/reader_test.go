package hiter_test

import (
	"encoding/hex"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/hiter/sh"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
)

func TestReader(t *testing.T) {
	src := slices.Collect(
		xiter.Limit(
			xiter.Map(
				func(b []byte) string { return hex.EncodeToString(b) },
				sh.RandBytes(8),
			),
			10,
		),
	)

	var (
		r   io.Reader
		err error
		bin []byte
	)

	r = hiter.Reader(
		func(s string) ([]byte, error) { return []byte(s), nil },
		iterable.NewResumable(slices.Values(src)).IntoIter(),
	)
	bin, err = io.ReadAll(r)
	assert.NilError(t, err)
	assert.Equal(t, strings.Join(src, ""), string(bin))

	r = hiter.Reader(
		func(s string) ([]byte, error) { return []byte(s + "\n"), nil },
		iterable.NewResumable(slices.Values(src)).IntoIter(),
	)
	bin, err = io.ReadAll(r)
	assert.NilError(t, err)
	assert.Equal(t, strings.Join(src, "\n")+"\n", string(bin))

	r = hiter.Reader(
		func(s string) ([]byte, error) { return []byte(s), nil },
		iterable.NewResumable(slices.Values(src)).IntoIter(),
	)
	err = iotest.TestReader(r, []byte(strings.Join(src, "")))
	assert.NilError(t, err)
}