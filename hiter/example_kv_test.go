package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Useful for test data.
func ExampleKeyValues() {
	kv := hiter.KeyValues[string, string]{{"foo", "bar"}, {"baz", "qux"}, {"quux", "corge"}}

	for k, v := range kv.Iter2() {
		fmt.Printf("%s:%s\n", k, v)
	}
	// Output:
	// foo:bar
	// baz:qux
	// quux:corge
}
