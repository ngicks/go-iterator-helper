package hiter_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestAppendBytes(t *testing.T) {
	src := [][]byte{[]byte(`bar`), []byte(`baz`)}
	assert.Assert(
		t,
		cmp.DeepEqual([]byte(`foobarbaz`), hiter.AppendBytes([]byte(`foo`), slices.Values(src))),
	)
}
