package hiter

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestCompact(t *testing.T) {
	type testCase struct {
		input    iter.Seq[int]
		expected []int
	}

	for _, tc := range []testCase{
		{
			input: Repeat(0, 0),
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
			input: Repeat(testSubject{}, 0),
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
		expected []KeyValue[int, int]
	}

	for _, tc := range []testCase{
		{
			input: Repeat2(0, 0, 0),
		},
		{
			input: Pairs( // 5,5 5,5 6,5 6,6 6,6
				xiter.Concat(Repeat(5, 2), Repeat(6, 3)),
				xiter.Concat(Repeat(5, 3), Repeat(6, 2)),
			),
			expected: KeyValues[int, int]{{K: 5, V: 5}, {K: 6, V: 5}, {K: 6, V: 6}},
		},
	} {
		assert.Assert(
			t,
			cmp.DeepEqual(
				tc.expected,
				Collect2(Compact2(tc.input)),
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
		expected []KeyValue[int, testSubject]
		eq       func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool
	}

	inputIter := slices.All([]testSubject{{"foo", 5}, {"foo", 6}, {"foo", 6}, {"bar", 6}, {"baz", 6}})

	for _, tc := range []testCase{
		{
			input: Repeat2(0, testSubject{}, 0),
		},
		{
			input: inputIter,
			expected: []KeyValue[int, testSubject]{
				{K: 0, V: testSubject{"foo", 5}},
				{K: 3, V: testSubject{"bar", 6}},
				{K: 4, V: testSubject{"baz", 6}},
			},
			eq: func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool { return v1.Foo == v2.Foo },
		},
		{
			input: inputIter,
			expected: []KeyValue[int, testSubject]{
				{K: 0, V: testSubject{"foo", 5}},
				{K: 1, V: testSubject{"foo", 6}},
			},
			eq: func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool { return v1.Bar == v2.Bar },
		},

		{
			input:    inputIter,
			expected: Collect2(inputIter),
			eq:       func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool { return k1 == k2 },
		},
	} {
		assert.Assert(t, cmp.DeepEqual(tc.expected, Collect2(CompactFunc2(tc.input, tc.eq))))
	}
}
