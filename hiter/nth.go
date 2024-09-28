package hiter

import "iter"

// Nth returns the nth element from seq by yielding up to n+1 elements.
// n is a 0-indexed number. Passing 0 as n returns the first element.
// Nth returns zero values if n is less than 0 or greater than or equals to the length of the iterator.
func Nth[V any](n int, seq iter.Seq[V]) (v V, ok bool) {
	if n < 0 {
		return
	}
	for v := range seq {
		if n == 0 {
			return v, true
		}
		n--
	}
	return
}

// Nth2 returns the nth pair from seq by yielding up to n+1 pairs.
// n is a 0-indexed number. Passing 0 as n returns the first element.
// Nth2 returns zero values if n is less than 0 or greater than or equals to the length of the iterator.
func Nth2[K, V any](n int, seq iter.Seq2[K, V]) (k K, v V, ok bool) {
	if n < 0 {
		return
	}
	for k, v := range seq {
		if n == 0 {
			return k, v, true
		}
		n--
	}
	return
}
