package iteratorhelper

import (
	"strconv"
	"testing"
)

func TestAdapter(t *testing.T) {
	// Chain
	testCases[int, int]{
		{
			iter:      Chain(SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}), SliceIter([]int{98, 2, 17, 5, 20, 8, 37})),
			expectedV: []int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3, 98, 2, 17, 5, 20, 8, 37},
		},
		{
			iter:      Chain(SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}), SliceIter([]int{98, 2, 17, 5, 20, 8, 37})),
			expectedV: []int{1, 5, 7, 8},
			breakIf:   func(k, v int) bool { return k > 3 },
		},
		{
			iter:      Chain(SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}), SliceIter([]int{98, 2, 17, 5, 20, 8, 37})),
			expectedV: []int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3, 98, 2, 17, 5},
			breakIf:   func(k, v int) bool { return v == 20 },
		},
	}.test(t)

	// Chunk
	testCases[[]int, []int]{
		{
			iter:      Chunk(SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}), 3),
			expectedV: [][]int{[]int{1, 5, 7}, []int{8, 1, 4}, []int{2, 6, 3}, []int{9, 3}},
		},
		{
			iter:      Chunk(SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}), 3),
			expectedV: [][]int{[]int{1, 5, 7}, []int{8, 1, 4}},
			breakIf:   func(k, v []int) bool { return v[0] == 2 },
		},
		{
			iter:      Chunk(SliceIter([]int{1, 5, 7}), 3),
			expectedV: [][]int{[]int{1, 5, 7}},
		},
		{
			iter:      Chunk(SliceIter([]int{1, 5, 7}), 4),
			expectedV: [][]int{[]int{1, 5, 7}},
		},
	}.test(t)

	// Enumerate
	testCases[int, int]{
		{
			iter:      Enumerate(SliceIterSingle([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3})),
			expectedK: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedV: []int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3},
		},
	}.test(t)

	// Filter
	testCases[int, int]{
		{
			iter: FilterSelect(
				SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}),
				func(k, v int) bool { return v%2 == 0 },
			),
			expectedV: []int{8, 4, 2, 6},
		},
		{
			iter: FilterExclude(
				SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}),
				func(k, v int) bool { return v%2 == 0 },
			),
			expectedV: []int{1, 5, 7, 1, 3, 9, 3},
		},
	}.test(t)

	// Map
	testCases[int, string]{
		{
			iter: Map(
				SliceIter([]int{1, 5, 7}),
				func(k, v int) (int, string) { return k, strconv.FormatInt(int64(v), 10) },
			),
			expectedV: []string{"1", "5", "7"},
		},
	}.test(t)

	// Skip/Take
	testCases[int, int]{
		{
			iter: SkipWhile(
				SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}),
				func(k, v int) bool { return v != 4 },
			),
			expectedV: []int{4, 2, 6, 3, 9, 3},
		},
		{
			iter: TakeWhile(
				SliceIter([]int{1, 5, 7, 8, 1, 4, 2, 6, 3, 9, 3}),
				func(k, v int) bool { return v != 4 },
			),
			expectedV: []int{1, 5, 7, 8, 1},
		},
	}.test(t)

	// Window
	testCases[[]int, []int]{
		{
			iter:      Window(SliceIter([]int{1, 5, 7, 8, 1, 4}), 3),
			expectedV: [][]int{[]int{1, 5, 7}, []int{5, 7, 8}, []int{7, 8, 1}, []int{8, 1, 4}},
		},
		{
			iter:      Window(SliceIter([]int{1, 5, 7, 8, 1, 4}), 3),
			expectedV: [][]int{[]int{1, 5, 7}, []int{5, 7, 8}, []int{7, 8, 1}},
			breakIf:   func(k, v []int) bool { return v[0] == 8 },
		},
		{
			iter:      Window(SliceIter([]int{1, 5}), 3),
			expectedV: [][]int{[]int{1, 5}},
		},
	}.test(t)

	// Zip
	testCases[int, string]{
		{
			iter:      Zip(SliceIterSingle([]int{1, 5, 7}), SliceIterSingle([]string{"foo", "bar", "baz"})),
			expectedK: []int{1, 5, 7},
			expectedV: []string{"foo", "bar", "baz"},
		},
		{
			iter:      Zip(SliceIterSingle([]int{1, 5, 7}), SliceIterSingle([]string{"foo", "bar", "baz", "qux"})),
			expectedK: []int{1, 5, 7},
			expectedV: []string{"foo", "bar", "baz"},
		},
		{
			iter:      Zip(SliceIterSingle([]int{1, 5, 7, 4}), SliceIterSingle([]string{"foo", "bar", "baz", "qux"})),
			expectedK: []int{1, 5},
			expectedV: []string{"foo", "bar"},
			breakIf:   func(l int, r string) bool { return r == "baz" },
		},
		{
			iter:      ZipPull(SliceIterSingle([]int{1, 5, 7}), SliceIterSingle([]string{"foo", "bar", "baz"})),
			expectedK: []int{1, 5, 7},
			expectedV: []string{"foo", "bar", "baz"},
		},
		{
			iter:      ZipPull(SliceIterSingle([]int{1, 5, 7}), SliceIterSingle([]string{"foo", "bar", "baz", "qux"})),
			expectedK: []int{1, 5, 7},
			expectedV: []string{"foo", "bar", "baz"},
		},
		{
			iter:      ZipPull(SliceIterSingle([]int{1, 5, 7, 4}), SliceIterSingle([]string{"foo", "bar", "baz", "qux"})),
			expectedK: []int{1, 5},
			expectedV: []string{"foo", "bar"},
			breakIf:   func(l int, r string) bool { return r == "baz" },
		},
	}.test(t)

	// ZipPair
	testCases[Pair[int, int], Pair[int, string]]{
		{
			iter:      ZipPair(SliceIter([]int{1, 5, 7}), SliceIter([]string{"foo", "bar", "baz"})),
			expectedK: []Pair[int, int]{Pair[int, int]{0, 1}, Pair[int, int]{1, 5}, Pair[int, int]{2, 7}},
			expectedV: []Pair[int, string]{Pair[int, string]{0, "foo"}, Pair[int, string]{1, "bar"}, Pair[int, string]{2, "baz"}},
		},
		{
			iter:      ZipPair(SliceIter([]int{1, 5, 7, 6}), SliceIter([]string{"foo", "bar", "baz"})),
			expectedK: []Pair[int, int]{Pair[int, int]{0, 1}, Pair[int, int]{1, 5}, Pair[int, int]{2, 7}},
			expectedV: []Pair[int, string]{Pair[int, string]{0, "foo"}, Pair[int, string]{1, "bar"}, Pair[int, string]{2, "baz"}},
		},
		{
			iter:      ZipPair(SliceIter([]int{1, 5, 7, 6}), SliceIter([]string{"foo", "bar", "baz"})),
			expectedK: []Pair[int, int]{Pair[int, int]{0, 1}, Pair[int, int]{1, 5}},
			expectedV: []Pair[int, string]{Pair[int, string]{0, "foo"}, Pair[int, string]{1, "bar"}},
			breakIf:   func(l Pair[int, int], r Pair[int, string]) bool { return r.V == "baz" },
		},
	}.test(t)

	// swap
	testCases[string, int]{
		{
			iter:      Swap(SliceIter([]string{"foo", "bar", "baz"})),
			expectedK: []string{"foo", "bar", "baz"},
			expectedV: []int{0, 1, 2},
		},
	}.test(t)

	// combined
	testCases[Pair[int, string], Pair[int, string]]{
		{
			iter: FilterExclude(
				ZipPair(
					Map[int, int, int, string](
						SliceIter([]int{1, 5, 7}),
						func(k, v int) (int, string) { return k, strconv.FormatInt(int64(v), 10) },
					),
					SliceIter([]string{"foo", "bar", "baz"}),
				),
				func(k, v Pair[int, string]) bool { return v.V == "bar" },
			),
			expectedK: []Pair[int, string]{Pair[int, string]{0, "1"}, Pair[int, string]{2, "7"}},
			expectedV: []Pair[int, string]{Pair[int, string]{0, "foo"}, Pair[int, string]{2, "baz"}},
		},
	}.test(t)
}
