package iterable

import (
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
