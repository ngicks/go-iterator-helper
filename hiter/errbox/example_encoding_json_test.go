package errbox_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing/iotest"

	"github.com/ngicks/go-iterator-helper/hiter/errbox"
)

// ExampleNewJsonDecoder_a_semantically_broken demonstrates raw decoder can be accessed while iterating over tokens.
// Also calling Decode is safe and not a race condition.
// Failing to decode does not affect its iteration.
// After the iterator stops, no error is stored.
func ExampleNewJsonDecoder_a_semantically_broken() {
	const semanticallyBroken = `{
		"foo": "bar",
		"baz": ["yay", "nay", 5, "wow"]
	}`

	dec := errbox.NewJsonDecoder(json.NewDecoder(strings.NewReader(semanticallyBroken)))
	defer dec.Stop()

	var depth int
	for t := range dec.IntoIter() {
		if depth == 1 && t == "baz" {
			// read opening [.
			t, err := dec.Dec.Token()
			if err != nil {
				panic(err)
			}
			if t != json.Delim('[') {
				panic("??")
			}
			var yayyay string
			for dec.Dec.More() {
				err = dec.Dec.Decode(&yayyay)
				if err == nil {
					fmt.Printf("yay? = %s\n", yayyay)
				} else {
					fmt.Printf("yay err = %v\n", err)
				}
			}
			// read closing ].
			t, err = dec.Dec.Token()
			if err != nil {
				panic(err)
			}
			if t != json.Delim(']') {
				panic("??")
			}
		}
		switch t {
		case json.Delim('{'), json.Delim('['):
			depth++
		case json.Delim('}'), json.Delim(']'):
			depth--
		}
	}
	fmt.Printf("stored error: %v\n", dec.Err())
	_, err := dec.Dec.Token()
	fmt.Printf("eof: %t\n", err == io.EOF)
	// Output:
	// yay? = yay
	// yay? = nay
	// yay err = json: cannot unmarshal number into Go value of type string
	// yay? = wow
	// stored error: <nil>
	// eof: true
}

// ExampleNewJsonDecoder_b_syntactically_broken demonstrates that
// syntactically broken json inputs cause no error on reading.
// Also it works well with some reader implementation where final non-empty data come with io.EOF error.
func ExampleNewJsonDecoder_b_syntactically_broken() {
	const syntacticallyBroken = `{
		"foo": "bar",
		"baz": ["yay", "nay", 5, "wow"],
		"broken": {
	}`

	dec := errbox.NewJsonDecoder(
		json.NewDecoder(
			iotest.DataErrReader(
				strings.NewReader(syntacticallyBroken),
			),
		),
	)
	defer dec.Stop()

	for t := range dec.IntoIter() {
		fmt.Printf("%v\n", t)
	}
	fmt.Printf("stored error: %v\n", dec.Err())
	_, err := dec.Dec.Token()
	fmt.Printf("eof: %t\n", err == io.EOF)
	// Output:
	// {
	// foo
	// bar
	// baz
	// [
	// yay
	// nay
	// 5
	// wow
	// ]
	// broken
	// {
	// }
	// stored error: <nil>
	// eof: true
}

// ExampleNewJsonDecoder_c_reader_broken demonstrates an error returned from the decoder
// can be inspected through Err method.
func ExampleNewJsonDecoder_c_reader_broken() {
	const readerBroken = `{
		"foo": "bar",
		"baz": ["yay", "nay", 5, "wow"]`

	dec := errbox.NewJsonDecoder(
		json.NewDecoder(
			io.MultiReader(
				strings.NewReader(readerBroken),
				iotest.ErrReader(errors.New("sample")),
			),
		),
	)
	defer dec.Stop()

	for t := range dec.IntoIter() {
		fmt.Printf("%v\n", t)
	}
	fmt.Printf("stored error: %v\n", dec.Err())
	// Output:
	// {
	// foo
	// bar
	// baz
	// [
	// yay
	// nay
	// 5
	// wow
	// ]
	// stored error: sample
}
