package reflectiter_test

import (
	"iter"
	"reflect"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/reflectiter"
)

type fields struct {
	A string
	B string
	C string
	D string
}

func TestFields(t *testing.T) {
	rv := reflect.ValueOf(fields{"A", "B", "C", "D"})
	rt := rv.Type()

	testcase.Two[reflect.StructField, reflect.Value]{
		Seq: func() iter.Seq2[reflect.StructField, reflect.Value] {
			return reflectiter.Fields(fields{"A", "B", "C", "D"})
		},
		Seqs: []func() iter.Seq2[reflect.StructField, reflect.Value]{
			func() iter.Seq2[reflect.StructField, reflect.Value] {
				return reflectiter.FieldsRv(reflect.ValueOf(fields{"A", "B", "C", "D"}))
			},
		},
		Expected: []hiter.KeyValue[reflect.StructField, reflect.Value]{
			{K: rt.Field(0), V: rv.Field(0)},
			{K: rt.Field(1), V: rv.Field(1)},
			{K: rt.Field(2), V: rv.Field(2)},
			{K: rt.Field(3), V: rv.Field(3)},
		},
		BreakAt: 1,
		CmpOpt:  []goCmp.Option{testcase.CompareReflectStructField, testcase.CompareReflectValue},
	}.Test(t)

	testcase.Two[reflect.StructField, reflect.Value]{
		Seq: func() iter.Seq2[reflect.StructField, reflect.Value] {
			return reflectiter.Fields(&fields{"A", "B", "C", "D"})
		},
		Seqs: []func() iter.Seq2[reflect.StructField, reflect.Value]{
			func() iter.Seq2[reflect.StructField, reflect.Value] {
				return reflectiter.FieldsRv(reflect.ValueOf(fields{"A", "B", "C", "D"}))
			},
		},
		Expected: []hiter.KeyValue[reflect.StructField, reflect.Value]{
			{K: rt.Field(0), V: rv.Field(0)},
			{K: rt.Field(1), V: rv.Field(1)},
			{K: rt.Field(2), V: rv.Field(2)},
			{K: rt.Field(3), V: rv.Field(3)},
		},
		BreakAt: 2,
		CmpOpt:  []goCmp.Option{testcase.CompareReflectStructField, testcase.CompareReflectValue},
	}.Test(t)
}
