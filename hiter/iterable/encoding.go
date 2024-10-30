package iterable

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ hiter.IntoIterable2[json.Token, error] = JsonDecoder{}
	_ hiter.IntoIterable2[xml.Token, error]  = XmlDecoder{}
)

// JsonDecoder adds IntoIter2 to [*json.Decoder].
type JsonDecoder struct {
	*json.Decoder
}

func (dec JsonDecoder) IntoIter2() iter.Seq2[json.Token, error] {
	return hiter.JsonDecoder(dec.Decoder)
}

// XmlDecoder adds IntoIter2 to [*xml.Decoder].
type XmlDecoder struct {
	*xml.Decoder
}

func (dec XmlDecoder) IntoIter2() iter.Seq2[xml.Token, error] {
	return hiter.XmlDecoder(dec.Decoder)
}

// CsvReader adds IntoIter2 method to [*csv.Reader].
type CsvReader struct {
	*csv.Reader
}

func (dec CsvReader) IntoIter2() iter.Seq2[[]string, error] {
	return hiter.CsvReader(dec.Reader)
}
