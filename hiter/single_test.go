package hiter_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestSingle(t *testing.T) {
	single := hiter.Single(5)
	single2 := hiter.Single2(3, 5)
	assert.Assert(t, cmp.DeepEqual([]int{5}, slices.Collect(single)))
	assert.Assert(t, cmp.DeepEqual([]hiter.KeyValue[int, int]{{K: 3, V: 5}}, hiter.Collect2(single2)))
}
