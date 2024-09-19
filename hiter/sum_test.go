package hiter

import (
	"slices"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSum(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	s := Sum(
		slices.Values(src),
	)
	assert.Equal(t, 15, s)
}

func TestSumOf(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	s := SumOf(
		func(ele int) int { return ele },
		slices.Values(src),
	)

	assert.Equal(t, 15, s)
}
