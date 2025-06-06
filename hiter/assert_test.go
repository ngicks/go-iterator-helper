package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
)

func TestAssert(t *testing.T) {
	src := []any{"foo", "bar", "baz"}
	src2 := []hiter.KeyValue[any, any]{{0, "foo"}, {1, "bar"}, {2, "baz"}}

	testcase.One[string]{
		Seq: func() iter.Seq[string] {
			return hiter.Assert[string](slices.Values(src))
		},
		Expected: []string{"foo", "bar", "baz"},
		BreakAt:  2,
	}.Test(t)

	testcase.Two[int, string]{
		Seq: func() iter.Seq2[int, string] {
			return hiter.Assert2[int, string](hiter.Values2(src2))
		},
		Expected: []hiter.KeyValue[int, string]{{0, "foo"}, {1, "bar"}, {2, "baz"}},
		BreakAt:  2,
	}.Test(t)
}
