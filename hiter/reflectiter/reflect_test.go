package reflectiter_test

import (
	"iter"
	"reflect"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/reflectiter"
)

func TestAssert(t *testing.T) {
	src := []any{"foo", "bar", "baz"}

	testcase.One[string]{
		Seq: func() iter.Seq[string] {
			rv := reflect.ValueOf(src)
			return reflectiter.AssertValue[string](hiter.OmitF(rv.Seq2()))
		},
		Expected: []string{"foo", "bar", "baz"},
		BreakAt:  2,
	}.Test(t)

	testcase.Two[int, string]{
		Seq: func() iter.Seq2[int, string] {
			rv := reflect.ValueOf(src)
			return reflectiter.AssertValue2[int, string](rv.Seq2())
		},
		Expected: []hiter.KeyValue[int, string]{{K: 0, V: "foo"}, {K: 1, V: "bar"}, {K: 2, V: "baz"}},
		BreakAt:  2,
	}.Test(t)
}
