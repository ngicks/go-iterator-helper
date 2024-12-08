package hiter_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestFind(t *testing.T) {
	src1 := []any{"foo", "bar", "foo", "baz", "foo"}
	src2 := []string{"foo", "bar", "foo", "baz", "foo"}
	src3 := [][]string{{"foo"}, {"foo", "bar"}, {"foo"}}

	type testCase struct {
		name string
		tgtV string
		idx  int
	}

	t.Run("Contains", func(t *testing.T) {
		assert.Assert(t, hiter.Contains("foo", slices.Values(src1)))
		assert.Assert(t, hiter.Contains("bar", slices.Values(src1)))
		assert.Assert(t, !hiter.Contains("qux", slices.Values(src1)))
	})

	t.Run("ContainsFunc", func(t *testing.T) {
		assert.Assert(t, hiter.ContainsFunc(func(s []string) bool { return slices.Contains(s, "foo") }, slices.Values(src3)))
		assert.Assert(t, !hiter.ContainsFunc(func(s []string) bool { return slices.Contains(s, "qux") }, slices.Values(src3)))
	})

	t.Run("Contains2", func(t *testing.T) {
		assert.Assert(t, hiter.Contains2(0, "foo", slices.All(src1)))
		assert.Assert(t, !hiter.Contains2(0, "bar", slices.All(src1)))
		assert.Assert(t, !hiter.Contains2(0, "qux", slices.All(src1)))
	})

	t.Run("ContainsFunc2", func(t *testing.T) {
		assert.Assert(t, hiter.ContainsFunc2(func(i int, s []string) bool { return i == 0 && slices.Contains(s, "foo") }, slices.All(src3)))
		assert.Assert(t, !hiter.ContainsFunc2(func(i int, s []string) bool { return i == 0 && slices.Contains(s, "bar") }, slices.All(src3)))
	})

	t.Run("Find", func(t *testing.T) {
		for _, tc := range []testCase{
			{
				"found_1",
				"foo",
				0,
			},
			{
				"found_2",
				"baz",
				3,
			},
			{
				"not_found",
				"yay",
				-1,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var (
					foundV any
					i      int
				)
				assertResult := func() {
					t.Helper()
					if i >= 0 {
						assert.Equal(t, tc.tgtV, foundV)
					} else {
						assert.Equal(t, nil, foundV)
					}
					assert.Equal(t, tc.idx, i)
				}
				foundV, i = hiter.Find(any(tc.tgtV), slices.Values(src1))
				assertResult()
				foundV, i = hiter.FindFunc(func(s any) bool { return s == tc.tgtV }, slices.Values(src1))
				assertResult()
			})
		}
	})

	t.Run("FindLast", func(t *testing.T) {
		for _, tc := range []testCase{
			{
				"found_1",
				"foo",
				4,
			},
			{
				"found_2",
				"baz",
				3,
			},
			{
				"not_found",
				"yay",
				-1,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var (
					foundV any
					i      int
				)
				assertResult := func() {
					t.Helper()
					if i >= 0 {
						assert.Equal(t, tc.tgtV, foundV)
					} else {
						assert.Equal(t, nil, foundV)
					}
					assert.Equal(t, tc.idx, i)
				}
				foundV, i = hiter.FindLast(any(tc.tgtV), slices.Values(src1))
				assertResult()
				foundV, i = hiter.FindLastFunc(func(s any) bool { return s == tc.tgtV }, slices.Values(src1))
				assertResult()
			})
		}
	})

	type testCase2 struct {
		name string
		tgtK int
		tgtV string
		idx  int
	}

	t.Run("Find2", func(t *testing.T) {
		for _, tc := range []testCase2{
			{
				"found_1",
				0,
				"foo",
				0,
			},
			{
				"found_1",
				1,
				"foo",
				-1,
			},
			{
				"found_2",
				3,
				"baz",
				3,
			},
			{
				"not_found",
				0,
				"yay",
				-1,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var (
					foundK int
					foundV string
					i      int
				)
				assertResult := func() {
					t.Helper()
					if i >= 0 {
						assert.Equal(t, tc.tgtK, foundK)
						assert.Equal(t, tc.tgtV, foundV)
					} else {
						assert.Equal(t, 0, foundK)
						assert.Equal(t, "", foundV)
					}
					assert.Equal(t, tc.idx, i)
				}
				foundK, foundV, i = hiter.Find2(tc.tgtK, tc.tgtV, slices.All(src2))
				assertResult()
				foundK, foundV, i = hiter.FindFunc2(func(i int, s string) bool { return i == tc.tgtK && s == tc.tgtV }, slices.All(src2))
				assertResult()
			})
		}
	})

	t.Run("FindLast2", func(t *testing.T) {
		for _, tc := range []testCase2{
			{
				"found_1",
				4,
				"foo",
				4,
			},
			{
				"found_1",
				1,
				"foo",
				-1,
			},
			{
				"found_2",
				3,
				"baz",
				3,
			},
			{
				"not_found",
				0,
				"yay",
				-1,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var (
					foundK int
					foundV string
					i      int
				)
				assertResult := func() {
					t.Helper()
					if i >= 0 {
						assert.Equal(t, tc.tgtK, foundK)
						assert.Equal(t, tc.tgtV, foundV)
					} else {
						assert.Equal(t, 0, foundK)
						assert.Equal(t, "", foundV)
					}
					assert.Equal(t, tc.idx, i)
				}
				foundK, foundV, i = hiter.FindLast2(tc.tgtK, tc.tgtV, slices.All(src2))
				assertResult()
				foundK, foundV, i = hiter.FindLastFunc2(func(i int, s string) bool { return i == tc.tgtK && s == tc.tgtV }, slices.All(src2))
				assertResult()
			})
		}
	})
}
