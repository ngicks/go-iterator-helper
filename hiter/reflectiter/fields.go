package reflectiter

import (
	"iter"
	"reflect"
)

// Fields is like [FieldsRv] but converts v to [reflect.Value].
func Fields(v any) iter.Seq2[reflect.StructField, reflect.Value] {
	return FieldsRv(reflect.ValueOf(v))
}

// FieldsRv returns an iterator over fields of struct value v.
// The iterator panics on invocation if v is not a struct or pointer to a struct value.
//
// Calling [reflect.Value.Interface] will panic if the field is not exported.
// Callers are advised to check if the field is exported by [reflect.StructField.IsExported].
func FieldsRv(rv reflect.Value) iter.Seq2[reflect.StructField, reflect.Value] {
	return func(yield func(reflect.StructField, reflect.Value) bool) {
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		rt := rv.Type()
		for i := range rv.NumField() {
			fty := rt.Field(i)
			fv := rv.Field(i)
			if !yield(fty, fv) {
				return
			}
		}
	}
}
