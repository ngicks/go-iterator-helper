package errbox

import (
	"encoding/json"
	"encoding/xml"

	"github.com/ngicks/go-iterator-helper/hiter"
)

type JsonDecoder struct {
	*Box[json.Token]
	Dec *json.Decoder
}

func NewJsonDecoder(dec *json.Decoder) *JsonDecoder {
	return &JsonDecoder{
		Box: New(hiter.JsonDecoder(dec)),
		Dec: dec,
	}
}

type XmlDecoder struct {
	*Box[xml.Token]
	Dec *xml.Decoder
}

func NewXmlDecoder(dec *xml.Decoder) *XmlDecoder {
	return &XmlDecoder{
		Box: New(hiter.XmlDecoder(dec)),
		Dec: dec,
	}
}
