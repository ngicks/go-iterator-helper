package hiter

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"io"
	"iter"
)

// JsonDecoder returns an iterator over json tokens.
// The first non-nil error encountered stops iteration after yielding it.
// [io.EOF] is excluded from result.
func JsonDecoder(dec *json.Decoder) iter.Seq2[json.Token, error] {
	return tokener(dec)
}

// XmlDecoder returns an iterator over xml tokens.
// The first non-nil error encountered stops iteration after yielding it.
// [io.EOF] is excluded from result.
// The caller should call [xml.CopyToken] before going to next iteration if they need to retain tokens.
func XmlDecoder(dec *xml.Decoder) iter.Seq2[xml.Token, error] {
	return tokener(dec)
}

func tokener[Dec interface{ Token() (V, error) }, V any](dec Dec) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		for {
			t, err := dec.Token()
			if err != nil {
				if err == io.EOF {
					return
				}
				yield(*new(V), err)
				return
			}
			if !yield(t, nil) {
				return
			}
		}
	}
}

// CsvReader returns an iterator over csv lines.
// Unlike [JsonDecoder] or [XmlDecoder], the iterator does not stop on first error
// since it could yield [csv.ErrFieldCount] with partial record and continue to next line.
//
// [io.EOF] is excluded from result.
func CsvReader(r *csv.Reader) iter.Seq2[[]string, error] {
	return reader(r)
}

func reader[R interface{ Read() (T, error) }, T any](r R) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for {
			rec, err := r.Read()
			if err == io.EOF {
				return
			}
			if !yield(rec, err) {
				return
			}
		}
	}
}
