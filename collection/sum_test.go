package collection

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
		slices.Values(src),
		func(ele int) int { return ele },
	)

	assert.Equal(t, 15, s)
}
