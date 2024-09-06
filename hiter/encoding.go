package hiter

import (
	"encoding/json"
	"encoding/xml"
	"iter"
)

// JsonDecoder returns an iterator over json tokens.
func JsonDecoder(dec *json.Decoder) iter.Seq2[json.Token, error] {
	return func(yield func(json.Token, error) bool) {
		for dec.More() {
			if !yield(dec.Token()) {
				return
			}
		}
		if !yield(dec.Token()) {
			return
		}
	}
}

// XmlDecoder returns an iterator over xml tokens.
// The first non-nil error encountered stops iteration.
// Callers should call [xml.CopyToken] before going to next iteration if they need to retain tokens.
func XmlDecoder(dec *xml.Decoder) iter.Seq2[xml.Token, error] {
	return func(yield func(xml.Token, error) bool) {
		for {
			tok, err := dec.Token()
			if !yield(tok, err) {
				return
			}
			if err != nil {
				return
			}
		}
	}
}
