package containeriter_test

import (
	"container/list"
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter/containeriter"
)

func ExampleListAll() {
	l := list.New()

	for _, s := range []string{"foo", "bar", "baz"} {
		l.PushBack(s)
	}

	fmt.Printf("all:              %#v\n", slices.Collect(containeriter.ListAll[string](l)))
	fmt.Printf("backward:         %#v\n", slices.Collect(containeriter.ListBackward[string](l)))
	fmt.Printf("element all:      %#v\n", slices.Collect(containeriter.ListElementAll[string](l.Front().Next())))
	fmt.Printf("element backward: %#v\n", slices.Collect(containeriter.ListElementBackward[string](l.Front().Next())))
	// Output:
	// all:              []string{"foo", "bar", "baz"}
	// backward:         []string{"baz", "bar", "foo"}
	// element all:      []string{"bar", "baz"}
	// element backward: []string{"bar", "foo"}
}
