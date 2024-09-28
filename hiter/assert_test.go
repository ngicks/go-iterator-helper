package hiter_test

import (
	"iter"
	"reflect"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func TestAssert(t *testing.T) {
	src := []any{"foo", "bar", "baz"}
	src2 := []hiter.KeyValue[any, any]{{0, "foo"}, {1, "bar"}, {2, "baz"}}

	testCase1[string]{
		Seq: func() iter.Seq[string] {
			return hiter.Assert[string](slices.Values(src))
		},
		Seqs: []func() iter.Seq[string]{
			func() iter.Seq[string] {
				rv := reflect.ValueOf(src)
				return hiter.AssertValue[string](hiter.OmitF(rv.Seq2()))
			},
		},
		Expected: []string{"foo", "bar", "baz"},
		BreakAt:  2,
	}.Test(t)

	testCase2[int, string]{
		Seq: func() iter.Seq2[int, string] {
			return hiter.Assert2[int, string](hiter.Values2(src2))
		},
		Seqs: []func() iter.Seq2[int, string]{
			func() iter.Seq2[int, string] {
				rv := reflect.ValueOf(src)
				return hiter.AssertValue2[int, string](rv.Seq2())
			},
		},
		Expected: []hiter.KeyValue[int, string]{{0, "foo"}, {1, "bar"}, {2, "baz"}},
		BreakAt:  2,
	}.Test(t)
}
