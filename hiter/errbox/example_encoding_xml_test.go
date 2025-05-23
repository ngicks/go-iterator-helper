package errbox_test

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing/iotest"

	"github.com/ngicks/go-iterator-helper/hiter/errbox"
)

// ExampleNewXmlDecoder_a_semantically_broken demonstrates raw decoder can be accessed while iterating over tokens.
// Also calling DecodeElement is safe and not a race condition.
// Failing to decode does not affect its iteration.
// After the iterator stops, no error is stored.
func ExampleNewXmlDecoder_a_semantically_broken() {
	const semanticallyBroken = `
	<root>
		<self/>
		<foo>bar</foo>
		<baz>5</baz>
		<baz>23</baz>
		<baz>yay</baz>
		<baz>49</baz>
	</root>`

	dec := errbox.NewXmlDecoder(xml.NewDecoder(strings.NewReader(strings.TrimSpace(semanticallyBroken))))

	var depth int
	for t := range dec.IntoIter() {
		var ok bool
		tok, ok := t.(xml.StartElement)
		if ok {
			if depth == 1 && tok.Name.Local == "baz" {
				var yayyay int
				err := dec.Dec.DecodeElement(&yayyay, &tok)
				if err == nil {
					fmt.Printf("yay? = %d\n", yayyay)
				} else {
					fmt.Printf("yay err = %v\n", err)
				}
				continue
			}
			depth++
		}
		_, ok = t.(xml.EndElement)
		if ok {
			depth--
		}
	}
	fmt.Printf("stored error: %v\n", dec.Err())
	_, err := dec.Dec.Token()
	fmt.Printf("eof: %t\n", err == io.EOF)
	// Output:
	// yay? = 5
	// yay? = 23
	// yay err = strconv.ParseInt: parsing "yay": invalid syntax
	// yay? = 49
	// stored error: <nil>
	// eof: true
}

// ExampleNewXmlDecoder_b_syntactically_broken demonstrates that
// syntactically broken xml inputs cause io.UnexpectedEOF error.
// Also it works well with some reader implementation where final non-empty data come with io.EOF error.
func ExampleNewXmlDecoder_b_syntactically_broken() {
	const syntacticallyBroken = `
	<root>
		<self/>
		<foo>bar</foo>
		<baz>5</baz>
		<baz>23</baz>
		<baz>yay</baz>
		<baz>49`

	dec := errbox.NewXmlDecoder(
		xml.NewDecoder(
			iotest.DataErrReader(
				strings.NewReader(strings.TrimSpace(syntacticallyBroken)),
			),
		),
	)

	for t := range dec.IntoIter() {
		fmt.Printf("%#v\n", t)
	}
	fmt.Printf("stored err: %v\n", dec.Err())
	// Output:
	// xml.StartElement{Name:xml.Name{Space:"", Local:"root"}, Attr:[]xml.Attr{}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"self"}, Attr:[]xml.Attr{}}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"self"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"foo"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x62, 0x61, 0x72}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"foo"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x35}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"baz"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x32, 0x33}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"baz"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x79, 0x61, 0x79}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"baz"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x34, 0x39}
	// stored err: XML syntax error on line 7: unexpected EOF
}

// ExampleNewXmlDecoder_c_reader_broken demonstrates an error returned from the decoder
// can be inspected through Err method.
func ExampleNewXmlDecoder_c_reader_broken() {
	const readerBroken = `
	<root>
		<self/>
		<foo>bar</foo>
		<baz>5</baz>
		<baz>23</baz>
		<baz>yay</baz>
		<baz>49`

	dec := errbox.NewXmlDecoder(
		xml.NewDecoder(
			io.MultiReader(
				strings.NewReader(strings.TrimSpace(readerBroken)),
				iotest.ErrReader(errors.New("sample")),
			),
		),
	)

	for t := range dec.IntoIter() {
		fmt.Printf("%#v\n", t)
	}
	fmt.Printf("stored err: %v\n", dec.Err())
	// Output:
	// xml.StartElement{Name:xml.Name{Space:"", Local:"root"}, Attr:[]xml.Attr{}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"self"}, Attr:[]xml.Attr{}}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"self"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"foo"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x62, 0x61, 0x72}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"foo"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x35}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"baz"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x32, 0x33}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"baz"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x79, 0x61, 0x79}
	// xml.EndElement{Name:xml.Name{Space:"", Local:"baz"}}
	// xml.CharData{0xa, 0x9, 0x9}
	// xml.StartElement{Name:xml.Name{Space:"", Local:"baz"}, Attr:[]xml.Attr{}}
	// xml.CharData{0x34, 0x39}
	// stored err: sample
}
