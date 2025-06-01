package hiter_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func ExampleFlatten() {
	yayyay := []string{"yay", "yay", "yay"}
	wah := []string{"wah", "wah"}
	first := true
	for s := range hiter.Flatten(hiter.Concat(slices.Values([][]string{yayyay}), slices.Values([][]string{wah}))) {
		if !first {
			fmt.Print(" ")
		}
		fmt.Print(s)
		first = false
	}
	// Output:
	// yay yay yay wah wah
}

func ExampleFlattenF() {
	yayyay := []string{"yay", "yay", "yay"}
	ohYah := []string{"oh", "yah"}
	first := true
	for k, v := range hiter.FlattenF(hiter.Pairs(hiter.Repeat(yayyay, -1), slices.Values(ohYah))) {
		if !first {
			fmt.Print(" ")
		}
		fmt.Printf("{%s %s}", k, v)
		first = false
	}
	// Output:
	// {yay oh} {yay oh} {yay oh} {yay yah} {yay yah} {yay yah}
}

func ExampleFlattenL() {
	yayyay := []string{"yay", "yay", "yay"}
	ohYah := []string{"oh", "yah"}
	first := true
	for k, v := range hiter.FlattenL(hiter.Pairs(slices.Values(yayyay), hiter.Repeat(ohYah, -1))) {
		if !first {
			fmt.Print(" ")
		}
		fmt.Printf("{%s %s}", k, v)
		first = false
	}
	// Output:
	// {yay oh} {yay yah} {yay oh} {yay yah} {yay oh} {yay yah}
}
