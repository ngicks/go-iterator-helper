package hiter

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"iter"
)

// JsonDecoder returns an iterator over json tokens.
func JsonDecoder(dec *json.Decoder) iter.Seq2[json.Token, error] {
	return tokener(dec)
}

// XmlDecoder returns an iterator over xml tokens.
// The first non-nil error encountered stops iteration.
// Callers should call [xml.CopyToken] before going to next iteration if they need to retain tokens.
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
