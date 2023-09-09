package iteratorhelper

import (
	"sort"
	"testing"

	cmp "github.com/google/go-cmp/cmp"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type testCase[K, V any] struct {
	iter      func(yield func(k K, v V) bool)
	expectedK []K
	expectedV []V
	breakIf   func(K, V) bool
	unordered bool
}

type testCases[K, V any] []testCase[K, V]

func (testCases testCases[K, V]) test(t *testing.T) {
	t.Helper()

	for _, tc := range testCases {
		var (
			outK []K
			outV []V
		)

		for k, v := range tc.iter {
			if tc.breakIf != nil && tc.breakIf(k, v) {
				break
			}
			outK = append(outK, k)
			outV = append(outV, v)
		}

		diffK := cmp.Diff(tc.expectedK, outK)
		diffV := cmp.Diff(tc.expectedV, outV)

		if tc.expectedK != nil && diffK != "" {
			t.Errorf("k not equal. diff = %s", diffK)
		}
		if tc.expectedV != nil && diffV != "" {
			t.Errorf("v not equal. diff = %s", diffV)
		}
	}
}

func TestBase(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	testCases[int, int]{
		{
			iter:      SliceIter([]int{1, 5, 7, 8}),
			expectedK: []int{0, 1, 2, 3},
			expectedV: []int{1, 5, 7, 8},
		},
		{
			iter:      SliceIter([]int{1, 5, 7, 8}),
			expectedK: []int{0, 1},
			expectedV: []int{1, 5},
			breakIf:   func(k, v int) bool { return k == 2 },
		},
		{
			iter:      Enumerate(SliceIterSingle([]int{1, 5, 7, 8})),
			expectedK: []int{0, 1, 2, 3},
			expectedV: []int{1, 5, 7, 8},
		},
		{
			iter:      Enumerate(SliceIterSingle([]int{1, 5, 7, 8})),
			expectedK: []int{0, 1},
			expectedV: []int{1, 5},
			breakIf:   func(k, v int) bool { return k == 2 },
		},
		{
			iter:      Enumerate(ChanIter(ch)),
			expectedV: []int{0, 1, 2, 3, 4},
		},
		{
			iter:      Enumerate(RangeIter(5, 9)),
			expectedV: []int{5, 6, 7, 8},
		},
	}.test(t)

	ordMap := orderedmap.New[string, string]()
	_ = ordMap.UnmarshalJSON([]byte(`{"foo":"foofoo","bar":"barbar","baz":"bazbaz"}`))

	testCases[string, string]{
		{
			iter:      OrderedMapIter[string, string](ordMap),
			expectedK: []string{"foo", "bar", "baz"},
			expectedV: []string{"foofoo", "barbar", "bazbaz"},
		},
	}.test(t)

	// Map
	iter := MapIter(map[int]int{0: 2, 1: 6, 2: 3})
	expectedK := []int{0, 1, 2}
	expectedV := []int{2, 3, 6}

	var outK, outV []int
	for k, v := range iter {
		outK = append(outK, k)
		outV = append(outV, v)
	}

	sort.Slice(outK, func(i, j int) bool { return outK[i] < outK[j] })
	sort.Slice(outV, func(i, j int) bool { return outV[i] < outV[j] })

	if diff := cmp.Diff(expectedK, outK); diff != "" {
		t.Errorf("not equal. diff =%s", diff)
	}
	if diff := cmp.Diff(expectedV, outV); diff != "" {
		t.Errorf("not equal. diff =%s", diff)
	}
}
