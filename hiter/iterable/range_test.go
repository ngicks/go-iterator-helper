package iterable_test

import (
	"math"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestRange_overflow(t *testing.T) {
	var r iterable.Range[int]

	r = iterable.Range[int]{Start: math.MaxInt, End: 0}
	_ = r.Reverse()

	r = iterable.Range[int]{Start: math.MinInt, End: 0}
	_ = r.Reverse()
}
