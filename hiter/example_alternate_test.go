package hiter_test

import (
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleAlternate() {
	tick := hiter.Repeat("tick", 3)
	tac := hiter.Repeat("tac", 3)

	for msg := range hiter.Alternate(tick, tac) {
		fmt.Printf("%s ", msg)
	}
	fmt.Printf("ooooohhhhh...\n")
	// Output:
	// tick tac tick tac tick tac ooooohhhhh...
}
