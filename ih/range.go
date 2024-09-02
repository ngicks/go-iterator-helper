package ih

import "iter"

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Range produces an iterator that yields sequential Numeric values in range [start, end).
// Values start from `start` and steps toward `end` 1 by 1,
// increased or decreased depending on start < end or not.
func Range[T Numeric](start, end T) iter.Seq[T] {
	return func(yield func(T) bool) {
		switch {
		default:
			return
		case start < end:
			for i := start; i < end; i++ {
				if !yield(i) {
					return
				}
			}
		case start > end:
			for i := start; i > end; i-- {
				if !yield(i) {
					return
				}
			}
		}
	}
}
