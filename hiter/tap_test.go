package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestTap(t *testing.T) {
	var observed1 []int

	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			observed1 = observed1[:0]
			return hiter.Tap(
				func(i int) {
					observed1 = append(observed1, i)
				},
				hiter.Range(0, 5),
			)
		},
		Expected: []int{0, 1, 2, 3, 4},
		BreakAt:  2,
	}.Test(t, func(_, count int) {
		switch count {
		case 0:
			assert.Assert(t, cmp.DeepEqual([]int{0, 1, 2, 3, 4}, observed1))
		case 1:
			assert.Assert(t, cmp.DeepEqual([]int{0, 1, 2}, observed1))
		}
	})

	var observed2 hiter.KeyValues[int, int]

	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			observed2 = observed2[:0]
			return hiter.Tap2(
				func(i, j int) {
					observed2 = append(observed2, hiter.KeyValue[int, int]{i, j})
				},
				hiter.Pairs(hiter.Range(5, 0), hiter.Range(0, 5)),
			)
		},
		Expected: hiter.KeyValues[int, int]{
			{5, 0},
			{4, 1},
			{3, 2},
			{2, 3},
			{1, 4},
		},
		BreakAt: 2,
	}.Test(t, func(_, count int) {
		switch count {
		case 0:
			assert.Assert(
				t,
				cmp.DeepEqual(
					hiter.KeyValues[int, int]{
						{5, 0},
						{4, 1},
						{3, 2},
						{2, 3},
						{1, 4},
					},
					observed2,
				),
			)
		case 1:
			assert.Assert(
				t,
				cmp.DeepEqual(
					hiter.KeyValues[int, int]{
						{5, 0},
						{4, 1},
						{3, 2},
					},
					observed2,
				),
			)
		}
	})
}

// TestTapLast tests TapLast function
func TestTapLast(t *testing.T) {
	t.Run("full iteration", func(t *testing.T) {
		var called bool
		seq := hiter.TapLast(
			func() { called = true },
			slices.Values([]int{1, 2, 3}),
		)

		result := slices.Collect(seq)

		assert.DeepEqual(t, []int{1, 2, 3}, result)
		assert.Assert(t, called, "TapLast should be called after full iteration")
	})

	t.Run("full iteration", func(t *testing.T) {
		var called bool
		seq := hiter.TapLast(
			func() { called = true },
			slices.Values([]int{1, 2, 3}),
		)

		result := slices.Collect(hiter.Limit(3, seq))
		assert.DeepEqual(t, []int{1, 2, 3}, result)
		assert.Assert(t, !called, "TapLast should not be called if loop was break-ed")
	})

	t.Run("break", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		type testCase struct {
			name string
			at   int
		}
		for _, tc := range []testCase{
			{
				"in middle",
				3,
			},
			{
				"at final",
				len(input),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var called bool
				seq := hiter.TapLast(
					func() { called = true },
					slices.Values(input),
				)

				result := slices.Collect(hiter.Limit(tc.at, seq))
				assert.DeepEqual(t, input[:tc.at], result)
				assert.Assert(t, !called, "TapLast should not be called if loop was break-ed")
			})
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		var called bool
		seq := hiter.TapLast(
			func() { called = true },
			slices.Values([]int{}),
		)

		result := slices.Collect(seq)

		assert.DeepEqual(t, []int(nil), result)
		assert.Assert(t, called, "TapLast should be called even for empty sequence")
	})
}

// TestTapLast2 tests TapLast2 function
func TestTapLast2(t *testing.T) {
	t.Run("full iteration", func(t *testing.T) {
		var called bool
		kvs := []hiter.KeyValue[string, int]{
			{K: "a", V: 1}, {K: "b", V: 2}, {K: "c", V: 3},
		}
		seq := hiter.TapLast2(
			func() { called = true },
			hiter.Values2(kvs),
		)

		result := hiter.Collect2(seq)

		assert.DeepEqual(t, kvs, result)
		assert.Assert(t, called, "TapLast2 should be called after full iteration")
	})

	t.Run("break", func(t *testing.T) {
		kvs := []hiter.KeyValue[string, int]{
			{K: "a", V: 1},
			{K: "b", V: 2},
			{K: "c", V: 3},
			{K: "d", V: 4},
		}
		type testCase struct {
			name string
			at   int
		}
		for _, tc := range []testCase{
			{
				"in middle",
				2,
			},
			{
				"at final",
				len(kvs),
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var called bool

				seq := hiter.TapLast2(
					func() { called = true },
					hiter.Values2(kvs),
				)

				result := hiter.Collect2(
					hiter.Limit2(tc.at, seq),
				)
				assert.DeepEqual(t, kvs[:tc.at], result)
				assert.Assert(t, !called, "TapLast2 should not be called if loop was breaked")
			})
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		var called bool
		seq := hiter.TapLast2(
			func() { called = true },
			hiter.Values2([]hiter.KeyValue[string, int]{}),
		)

		result := hiter.Collect2(seq)

		assert.DeepEqual(t, []hiter.KeyValue[string, int](nil), result)
		assert.Assert(t, called, "TapLast2 should be called even for empty sequence")
	})
}
