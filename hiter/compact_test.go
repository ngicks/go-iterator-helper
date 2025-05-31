package hiter

import (
	"iter"
	"slices"
	"testing"

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

	// Test early termination
	t.Run("early termination", func(t *testing.T) {
		var collected []int
		seq := Compact(slices.Values([]int{1, 1, 2, 2, 3, 3}))
		for v := range seq {
			collected = append(collected, v)
			if len(collected) >= 2 {
				break
			}
		}
		assert.DeepEqual(t, []int{1, 2}, collected)
	})
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
		assert.Assert(t, cmp.DeepEqual(tc.expected, slices.Collect(CompactFunc(tc.eq, tc.input))))
	}

	// Test early termination
	t.Run("early termination", func(t *testing.T) {
		var collected []testSubject
		eq := func(i, j testSubject) bool { return i.Foo == j.Foo }
		seq := CompactFunc(eq, slices.Values([]testSubject{{"foo", 5}, {"foo", 6}, {"bar", 6}, {"baz", 6}}))
		for v := range seq {
			collected = append(collected, v)
			if len(collected) >= 2 {
				break
			}
		}
		assert.DeepEqual(t, []testSubject{{"foo", 5}, {"bar", 6}}, collected)
	})
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
				Concat(Repeat(5, 2), Repeat(6, 3)),
				Concat(Repeat(5, 3), Repeat(6, 2)),
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

	// Test early termination
	t.Run("early termination", func(t *testing.T) {
		var collected []KeyValue[int, int]
		seq := Compact2(Pairs(
			slices.Values([]int{1, 1, 2, 2, 3, 3}),
			slices.Values([]int{1, 1, 2, 2, 3, 3}),
		))
		for k, v := range seq {
			collected = append(collected, KeyValue[int, int]{K: k, V: v})
			if len(collected) >= 2 {
				break
			}
		}
		assert.DeepEqual(t, []KeyValue[int, int]{{K: 1, V: 1}, {K: 2, V: 2}}, collected)
	})
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
		assert.Assert(t, cmp.DeepEqual(tc.expected, Collect2(CompactFunc2(tc.eq, tc.input))))
	}

	// Test early termination
	t.Run("early termination", func(t *testing.T) {
		var collected []KeyValue[int, testSubject]
		eq := func(k1 int, v1 testSubject, k2 int, v2 testSubject) bool { return v1.Foo == v2.Foo }
		seq := CompactFunc2(eq, slices.All([]testSubject{{"foo", 5}, {"foo", 6}, {"bar", 6}, {"baz", 6}}))
		for k, v := range seq {
			collected = append(collected, KeyValue[int, testSubject]{K: k, V: v})
			if len(collected) >= 2 {
				break
			}
		}
		expected := []KeyValue[int, testSubject]{
			{K: 0, V: testSubject{"foo", 5}},
			{K: 2, V: testSubject{"bar", 6}},
		}
		assert.DeepEqual(t, expected, collected)
	})
}
