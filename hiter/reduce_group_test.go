package hiter

import (
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestReduceGroup(t *testing.T) {
	kvs := KeyValues[string, []int]{
		{K: "foo", V: []int{1, 2, 3}},
		{K: "bar", V: []int{1, 2, 3}},
		{K: "baz", V: []int{1, 2, 3}},
		{K: "foo", V: []int{3, 4, 5}},
		{K: "bar", V: []int{1, 5, 5}},
		{K: "foo", V: []int{12}},
	}
	m := ReduceGroup(
		func(accum []int, c []int) []int {
			return append(accum, c...)
		},
		[]int{0},
		kvs.Iter2(),
	)
	assert.Assert(
		t,
		cmp.DeepEqual(
			map[string][]int{
				"foo": {0, 1, 2, 3, 3, 4, 5, 12},
				"bar": {0, 1, 2, 3, 1, 5, 5},
				"baz": {0, 1, 2, 3},
			},
			m,
		),
	)

	m = map[string][]int{"wow": {5}}
	m = InsertReduceGroup(
		m,
		func(accum []int, c []int) []int {
			return append(accum, c...)
		},
		[]int{0},
		kvs.Iter2(),
	)
	assert.Assert(
		t,
		cmp.DeepEqual(
			map[string][]int{
				"foo": {0, 1, 2, 3, 3, 4, 5, 12},
				"bar": {0, 1, 2, 3, 1, 5, 5},
				"baz": {0, 1, 2, 3},
				"wow": {5},
			},
			m,
		),
	)

}
