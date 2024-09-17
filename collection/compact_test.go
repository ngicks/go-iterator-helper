package collection

import (
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func Example_compact() {
	m := xiter.Merge(
		xiter.Map(func(i int) int { return 2 * i }, hiter.Range(1, 11)),
		xiter.Map(func(i int) int { return 1 << i }, hiter.Range(1, 11)),
	)

	first := true
	for i := range Compact(m) {
		if !first {
			fmt.Printf(", ")
		}
		fmt.Printf("%d", i)
		first = false
	}
	fmt.Println()
	// Output:
	// 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 32, 64, 128, 256, 512, 1024
}

func TestCompact(t *testing.T) {
	type testCase struct {
		input    iter.Seq[int]
		expected []int
	}

	for _, tc := range []testCase{
		{
			input: hiter.Repeat(0, 0),
		},
		{
			input:    slices.Values([]int{1, 1, 2, 2, 5, 5, 2, 2, 3, 3, 4}),
			expected: []int{1, 2, 5, 2, 3, 4},
		},
	} {
		assert.Assert(t, cmp.DeepEqual(tc.expected, slices.Collect(Compact(tc.input))))
	}
}

func TestCompactFunc(t *testing.T) {
	type testSubject struct {
		Foo string
		Bar int
	}

	type testCase struct {
		input    iter.Seq[testSubject]
		expected []testSubject
		eq       func(i, j testSubject) bool
	}

	for _, tc := range []testCase{
		{
			input: hiter.Repeat(testSubject{}, 0),
		},
		{
			input:    slices.Values([]testSubject{{"foo", 5}, {"foo", 6}, {"foo", 6}, {"bar", 6}, {"baz", 6}}),
			expected: []testSubject{{"foo", 5}, {"bar", 6}, {"baz", 6}},
			eq:       func(i, j testSubject) bool { return i.Foo == j.Foo },
		},
		{
			input:    slices.Values([]testSubject{{"foo", 5}, {"foo", 6}, {"foo", 6}, {"bar", 6}, {"baz", 6}}),
			expected: []testSubject{{"foo", 5}, {"foo", 6}},
			eq:       func(i, j testSubject) bool { return i.Bar == j.Bar },
		},
	} {
		assert.Assert(t, cmp.DeepEqual(tc.expected, slices.Collect(CompactFunc(tc.input, tc.eq))))
	}
}

func TestCompact2(t *testing.T) {
	type testCase struct {
		input    iter.Seq2[int, int]
		expected []hiter.KeyValue[int, int]
	}

	for _, tc := range []testCase{
		{
			input: hiter.Repeat2(0, 0, 0),
		},
		{
			input: hiter.Pairs( // 5,5 5,5 6,5 6,6 6,6
				xiter.Concat(hiter.Repeat(5, 2), hiter.Repeat(6, 3)),
				xiter.Concat(hiter.Repeat(5, 3), hiter.Repeat(6, 2)),
			),
			expected: hiter.KeyValues[int, int]{{K: 5, V: 5}, {K: 6, V: 5}, {K: 6, V: 6}},
		},
	} {
		assert.Assert(
			t,
			cmp.DeepEqual(
				tc.expected,
				hiter.Collect2(Compact2(tc.input)),
			),
		)
	}
}

func TestCompactFunc2(t *testing.T) {
	type testSubject struct {
		Foo string
		Bar int
	}

	type testCase struct {
		input    iter.Seq2[int, testSubject]
		expected []hiter.KeyValue[int, testSubject]
		eq       func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool
	}

	inputIter := slices.All([]testSubject{{"foo", 5}, {"foo", 6}, {"foo", 6}, {"bar", 6}, {"baz", 6}})

	for _, tc := range []testCase{
		{
			input: hiter.Repeat2(0, testSubject{}, 0),
		},
		{
			input: inputIter,
			expected: []hiter.KeyValue[int, testSubject]{
				{K: 0, V: testSubject{"foo", 5}},
				{K: 3, V: testSubject{"bar", 6}},
				{K: 4, V: testSubject{"baz", 6}},
			},
			eq: func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool { return v1.Foo == v2.Foo },
		},
		{
			input: inputIter,
			expected: []hiter.KeyValue[int, testSubject]{
				{K: 0, V: testSubject{"foo", 5}},
				{K: 1, V: testSubject{"foo", 6}},
			},
			eq: func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool { return v1.Bar == v2.Bar },
		},

		{
			input:    inputIter,
			expected: hiter.Collect2(inputIter),
			eq:       func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool { return k1 == k2 },
		},
	} {
		assert.Assert(t, cmp.DeepEqual(tc.expected, hiter.Collect2(CompactFunc2(tc.input, tc.eq))))
	}
}
