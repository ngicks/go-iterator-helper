package hiter_test

import (
	"encoding/json"
	"encoding/xml"
	"iter"
	"strings"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
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
		var box *errbox.JsonDecoder
		testCase2[json.Token, error]{
			Seq: func() iter.Seq2[json.Token, error] {
				return hiter.JsonDecoder(jsonDec())
			},
			Seqs: []func() iter.Seq2[json.Token, error]{
				func() iter.Seq2[json.Token, error] {
					return iterable.JsonDecoder{Decoder: jsonDec()}.IntoIter2()
				},
				func() iter.Seq2[json.Token, error] {
					box = errbox.NewJsonDecoder(jsonDec())
					return hiter.Pairs(box.IntoIter(), hiter.Repeat(error(nil), -1))
				},
			},
			Expected: []hiter.KeyValue[json.Token, error]{
				{json.Delim('{'), nil},
				{json.Token("foo"), nil}, {json.Token("bar"), nil},
				{json.Delim('}'), nil},
			},
			BreakAt:  2,
			Stateful: true,
		}.Test(t, func(_, _ int) {
			if box != nil {
				assert.NilError(t, box.Err())
			}
		})
	})

	t.Run("XML", func(t *testing.T) {
		var box *errbox.XmlDecoder
		testCase2[xml.Token, error]{
			Seq: func() iter.Seq2[xml.Token, error] {
				return xiter.Map2(copyXmlToken, hiter.XmlDecoder(xmlDec()))
			},
			Seqs: []func() iter.Seq2[xml.Token, error]{
				func() iter.Seq2[xml.Token, error] {
					return xiter.Map2(copyXmlToken, iterable.XmlDecoder{Decoder: xmlDec()}.IntoIter2())
				},
				func() iter.Seq2[xml.Token, error] {
					box = errbox.NewXmlDecoder(xmlDec())
					return hiter.Pairs(xiter.Map(xml.CopyToken, box.IntoIter()), hiter.Repeat(error(nil), -1))
				},
			},
			Expected: []hiter.KeyValue[xml.Token, error]{
				{xml.StartElement{Name: xml.Name{Local: "foo"}, Attr: []xml.Attr{}}, nil},
				{xml.CharData("yay"), nil},
				{xml.EndElement{Name: xml.Name{Local: "foo"}}, nil},
				{xml.StartElement{Name: xml.Name{Local: "bar"}, Attr: []xml.Attr{}}, nil},
				{xml.CharData("nay"), nil},
				{xml.EndElement{Name: xml.Name{Local: "bar"}}, nil},
			},
			BreakAt:  2,
			CmpOpt:   []goCmp.Option{compareErrorsIs},
			Stateful: true,
		}.Test(t)
	})
}
