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

func Example_encoding_json_semantically_broken() {
	const semanticallyBroken = `{
		"foo": "bar",
		"baz": ["yay", "nay", 5, "wow"]
	}`

	dec := errbox.NewJsonDecoder(json.NewDecoder(strings.NewReader(semanticallyBroken)))

	var depth int
	for t := range dec.Iter() {
		if depth == 1 && t == "baz" {
			// read 1 ahead.
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

func Example_encoding_json_syntactically_broken() {
	const syntacticallyBroken = `{
		"foo": "bar",
		"baz": ["yay", "nay", 5, "wow"],
		"broken": {
	}`

	dec := errbox.NewJsonDecoder(json.NewDecoder(strings.NewReader(syntacticallyBroken)))

	for t := range dec.Iter() {
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

func Example_encoding_json_reader_broken() {
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

	for t := range dec.Iter() {
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
