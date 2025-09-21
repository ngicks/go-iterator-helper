package goiteratorhelper_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter/encodingiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
)

// Example_errboxJSON demonstrates error handling in JSON streams
func Example_errboxJSON() {
	// JSON stream with an invalid entry
	jsonData := `{"id":1,"name":"Alice"}
{"id":2,"name":"Bob"}
{"invalid json syntax
{"id":3,"name":"Charlie"}`

	type Person struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonData)))
	jsonBox := errbox.New(encodingiter.Decode[Person](decoder))

	var validRecords []Person
	for person := range jsonBox.IntoIter() {
		validRecords = append(validRecords, person)
		fmt.Printf("Decoded: %+v\n", person)
	}

	// Check for errors after processing
	if err := jsonBox.Err(); err != nil {
		fmt.Printf("Stream error after %d valid records: %v\n", len(validRecords), err)
	}

	// Output:
	// Decoded: {ID:1 Name:Alice}
	// Decoded: {ID:2 Name:Bob}
	// Stream error after 2 valid records: invalid character '\n' in string literal
}