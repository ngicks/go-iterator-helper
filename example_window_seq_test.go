package goiteratorhelper_test

import (
	"fmt"
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_windowSeq demonstrates WindowSeq with iterator-based windows
func Example_windowSeq() {
	// Temperature readings
	temperatures := []float64{22.5, 23.0, 22.8, 23.5, 24.0, 23.8, 23.2}

	// Create iterator-based sliding windows
	windowSeqs := hiter.WindowSeq(3, slices.Values(temperatures))

	// Find max temperature in each window
	maxTemps := hiter.Map(func(window iter.Seq[float64]) float64 {
		max := 0.0
		for temp := range window {
			if temp > max {
				max = temp
			}
		}
		return max
	}, windowSeqs)

	fmt.Printf("Max temperatures in sliding windows: %.1f\n", slices.Collect(maxTemps))

	// Output:
	// Max temperatures in sliding windows: [23.0 23.5 24.0 24.0 24.0]
}