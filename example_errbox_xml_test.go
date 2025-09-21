package goiteratorhelper_test

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter/encodingiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
)

// Example_errboxXML demonstrates error handling in XML streams
func Example_errboxXML() {
	// XML stream with mixed valid and invalid entries
	xmlData := `<item><id>1</id><name>First</name></item>
<item><id>2</id><name>Second</name></item>
<item><id>3<name>Third</name></item>
<item><id>4</id><name>Fourth</name></item>`

	type Item struct {
		ID   int    `xml:"id"`
		Name string `xml:"name"`
	}

	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlData)))
	xmlBox := errbox.New(encodingiter.Decode[Item](decoder))

	validCount := 0
	for item := range xmlBox.IntoIter() {
		validCount++
		fmt.Printf("Valid item: %+v\n", item)
	}

	if err := xmlBox.Err(); err != nil {
		fmt.Printf("Parsing stopped after %d valid items: error in XML\n", validCount)
	}

	// Output:
	// Valid item: {ID:1 Name:First}
	// Valid item: {ID:2 Name:Second}
	// Parsing stopped after 2 valid items: error in XML
}