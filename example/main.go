package main

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/iteratorhelper"
)

func main() {
	var iter func(yield func(k int, v string) bool)
	iter = iteratorhelper.SliceIter([]string{"foo", "bar", "baz", "qux", "quux"})
	var adaptedIter func(yield func(k []int, v []string) bool)
	adaptedIter = iteratorhelper.Chunk(iter, 2)

	for k, v := range adaptedIter {
		fmt.Printf("k = %#v, v = %#v\n", k, v)
	}
}
