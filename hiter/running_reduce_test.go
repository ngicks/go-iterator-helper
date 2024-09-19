package hiter

import (
	"slices"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestRunningReduce(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	r := slices.Collect(
		RunningReduce(
			func(accum int, next int, idx int) int { return accum + next },
			0,
			slices.Values(src),
		),
	)
	assert.Assert(
		t,
		cmp.DeepEqual(
			[]int{1, 3, 6, 10, 15},
			r,
		),
	)
}
