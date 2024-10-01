package hiter

import "iter"

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Range produces an iterator that yields sequential Numeric values in range [start, end).
// Values start from `start` and step toward `end`.
// At each step value is increased by 1 if start < end, otherwise decreased by 1.
func Range[T Numeric](start, end T) iter.Seq[T] {
	return rangeInclusive(start, end, true, false)
}

func RangeInclusive[T Numeric](start, end T, includeStart, includeEnd bool) iter.Seq[T] {
	return rangeInclusive(start, end, includeStart, includeEnd)
}

func rangeInclusive[T Numeric](start, end T, includeStart, includeEnd bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		start := start
		end := end
		switch {
		default:
			return
		case start < end:
			if !includeStart {
				start += 1
			}
			if !includeEnd {
				end -= 1
			}
			for i := start; i <= end; i++ {
				if !yield(i) {
					return
				}
			}
		case start > end:
			if !includeStart {
				start -= 1
			}
			if !includeEnd {
				end += 1
			}
			for i := start; i >= end; i-- {
				if !yield(i) {
					return
				}
			}
		}
	}
}
