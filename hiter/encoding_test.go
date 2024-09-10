package hiter_test

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"iter"
	"strings"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

func TestEncoding(t *testing.T) {
	jsonDec := func() *json.Decoder {
		return json.NewDecoder(strings.NewReader(`{"foo":"bar"}`))
	}
	xmlDec := func() *xml.Decoder {
		return xml.NewDecoder(strings.NewReader(`<foo>yay</foo><bar>nay</bar>`))
	}
	copyXmlToken := func(x xml.Token, e error) (xml.Token, error) {
		return xml.CopyToken(x), e
	}
	t.Run("JSON", func(t *testing.T) {
		testCase2[json.Token, error]{
			Seq: func() iter.Seq2[json.Token, error] {
				return hiter.JsonDecoder(jsonDec())
			},
			Seqs: []func() iter.Seq2[json.Token, error]{
				func() iter.Seq2[json.Token, error] {
					return iterable.JsonDecoder{Decoder: jsonDec()}.IntoIter2()
				},
			},
			Expected: []hiter.KeyValue[json.Token, error]{
				{json.Delim('{'), nil},
				{json.Token("foo"), nil}, {json.Token("bar"), nil},
				{json.Delim('}'), nil},
			},
			BreakAt:  2,
			Stateful: true,
		}.Test(t)
	})

	t.Run("XML", func(t *testing.T) {
		testCase2[xml.Token, error]{
			Seq: func() iter.Seq2[xml.Token, error] {
				return xiter.Map2(copyXmlToken, hiter.XmlDecoder(xmlDec()))
			},
			Seqs: []func() iter.Seq2[xml.Token, error]{
				func() iter.Seq2[xml.Token, error] {
					return xiter.Map2(copyXmlToken, iterable.XmlDecoder{Decoder: xmlDec()}.IntoIter2())
				},
			},
			Expected: []hiter.KeyValue[xml.Token, error]{
				{xml.StartElement{Name: xml.Name{Local: "foo"}, Attr: []xml.Attr{}}, nil},
				{xml.CharData("yay"), nil},
				{xml.EndElement{Name: xml.Name{Local: "foo"}}, nil},
				{xml.StartElement{Name: xml.Name{Local: "bar"}, Attr: []xml.Attr{}}, nil},
				{xml.CharData("nay"), nil},
				{xml.EndElement{Name: xml.Name{Local: "bar"}}, nil},
				{nil, io.EOF},
			},
			BreakAt:  2,
			CmpOpt:   []goCmp.Option{goCmp.Comparer(func(e1, e2 error) bool { return errors.Is(e1, e2) })},
			Stateful: true,
		}.Test(t)
	})

}