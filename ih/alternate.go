package ih

import "iter"

// Alternate returns an iterator that yields alternatively each seq from head to tail.
// The first exhausted seq stops the iterator.
func Alternate[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		if len(seqs) == 0 {
			return
		}

		nexts := make([]func() (V, bool), len(seqs)-1)
		stops := make([]func(), len(seqs)-1)

		for i, it := range seqs[1:] {
			next, stop := iter.Pull(it)
			nexts[i] = next
			stops[i] = stop
		}
		defer func() {
			for _, s := range stops {
				s()
			}
		}()

		for v := range seqs[0] {
			if !yield(v) {
				return
			}
			for _, n := range nexts {
				v, ok := n()
				if !ok {
					return
				}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Alternate returns an iterator that yields alternatively each seq from head to tail.
// The first exhausted seq stops the iterator.
func Alternate2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if len(seqs) == 0 {
			return
		}

		nexts := make([]func() (K, V, bool), len(seqs)-1)
		stops := make([]func(), len(seqs)-1)

		for i, seq := range seqs[1:] {
			next, stop := iter.Pull2(seq)
			nexts[i] = next
			stops[i] = stop
		}
		defer func() {
			for _, s := range stops {
				s()
			}
		}()

		for k, v := range seqs[0] {
			if !yield(k, v) {
				return
			}
			for _, n := range nexts {
				k, v, ok := n()
				if !ok {
					return
				}
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
