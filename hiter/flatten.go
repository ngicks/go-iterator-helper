package hiter

import "iter"

// Flatten returns an iterator over each value of slices from seq.
func Flatten[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for s := range seq {
			for _, t := range s {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// Flatten returns and iterator over values from iterators from seq.
func FlattenSeq[V any](seq iter.Seq[iter.Seq[V]]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for s := range seq {
			for t := range s {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// FlattenSeq2 returns and iterator over key-value pairs from iterators from seq.
func FlattenSeq2[K, V any](seq iter.Seq[iter.Seq2[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for s := range seq {
			for k, v := range s {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// FlattenF returns an iterator over flattened key-value pairs from seq.
// While the iterator iterates over the former value from seq, the latter value is repeated.
func FlattenF[S1 ~[]E1, E1 any, E2 any](seq iter.Seq2[S1, E2]) iter.Seq2[E1, E2] {
	return func(yield func(E1, E2) bool) {
		for s, g := range seq {
			for _, e := range s {
				if !yield(e, g) {
					return
				}
			}
		}
	}
}

// FlattenL returns an iterator over flattened key-value pairs from seq.
// While the iterator iterates over the latter value from seq, the former value is repeated.
func FlattenL[S2 ~[]E2, E1 any, E2 any](seq iter.Seq2[E1, S2]) iter.Seq2[E1, E2] {
	return func(yield func(E1, E2) bool) {
		for e, t := range seq {
			for _, g := range t {
				if !yield(e, g) {
					return
				}
			}
		}
	}
}

// FlattenSeqF returns an iterator over flattened key-value pairs from seq.
// While the iterator iterates over the former value from seq, the latter value is repeated.
func FlattenSeqF[K, V any](seq iter.Seq2[iter.Seq[K], V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			for k := range k {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// FlattenSeqF returns an iterator over flattened pairs from seq.
// While the iterator iterates over the latter value from seq, the former value is repeated.
func FlattenSeqL[K, V any](seq iter.Seq2[K, iter.Seq[V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			for v := range v {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
