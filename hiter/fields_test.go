package hiter_test

import (
	"iter"
	"reflect"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
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
			return hiter.Fields(fields{"A", "B", "C", "D"})
		},
		Expected: []hiter.KeyValue[reflect.StructField, reflect.Value]{
			{rt.Field(0), rv.Field(0)},
			{rt.Field(1), rv.Field(1)},
			{rt.Field(2), rv.Field(2)},
			{rt.Field(3), rv.Field(3)},
		},
		BreakAt: 1,
		CmpOpt:  []goCmp.Option{testcase.CompareReflectStructField, testcase.CompareReflectValue},
	}.Test(t)

	testcase.Two[reflect.StructField, reflect.Value]{
		Seq: func() iter.Seq2[reflect.StructField, reflect.Value] {
			return hiter.Fields(&fields{"A", "B", "C", "D"})
		},
		Expected: []hiter.KeyValue[reflect.StructField, reflect.Value]{
			{rt.Field(0), rv.Field(0)},
			{rt.Field(1), rv.Field(1)},
			{rt.Field(2), rv.Field(2)},
			{rt.Field(3), rv.Field(3)},
		},
		BreakAt: 2,
		CmpOpt:  []goCmp.Option{testcase.CompareReflectStructField, testcase.CompareReflectValue},
	}.Test(t)
}
