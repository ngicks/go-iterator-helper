package hiter_test

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleReduceGroup() {
	m1 := map[int]string{
		0: "foo",
		1: "bar",
		2: "baz",
	}
	m2 := map[int]string{
		0: "foo",
		2: "zab",
		3: "gooood",
	}

	reduced := hiter.ReduceGroup(
		func(sum []string, c string) []string { return append(sum, c) },
		nil,
		hiter.Concat2(maps.All(m1), maps.All(m2)),
	)

	bin, _ := json.MarshalIndent(reduced, "", "    ")

	fmt.Printf("%s\n", bin)
	// Output:
	// {
	//     "0": [
	//         "foo",
	//         "foo"
	//     ],
	//     "1": [
	//         "bar"
	//     ],
	//     "2": [
	//         "baz",
	//         "zab"
	//     ],
	//     "3": [
	//         "gooood"
	//     ]
	// }
}
