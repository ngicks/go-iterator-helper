package goiteratorhelper_test

import (
	"fmt"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Example_windowMovingAverage demonstrates sliding window for moving averages
func Example_windowMovingAverage() {
	// Data for moving average calculation
	data := []int{10, 20, 30, 40, 50, 60, 70, 80}
	windowSize := 3

	// Create sliding windows
	windows := hiter.Window(data, windowSize)

	// Calculate average for each window
	averages := hiter.Map(func(window []int) float64 {
		sum := hiter.Sum(slices.Values(window))
		return float64(sum) / float64(len(window))
	}, windows)

	fmt.Println("Data:", data)
	fmt.Printf("Moving averages (window=%d): %.1f\n", windowSize, slices.Collect(averages))

	// Output:
	// Data: [10 20 30 40 50 60 70 80]
	// Moving averages (window=3): [20.0 30.0 40.0 50.0 60.0 70.0]
}