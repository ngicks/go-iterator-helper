package hiter

import "iter"

// TODO: consider adding
//  - https://deno.land/std@0.224.0/collections/associate_by.ts?s=associateBy
//  - https://deno.land/std@0.224.0/collections/associate_with.ts?s=associateWith
//

// Permutations returns an iterator that yields permutations of in.
// The returned iterator reorders in in-place.
// The caller should not retain in or slices from the iterator,
// Or should explicitly clone yielded values.
func Permutations[S ~[]E, E any](in S) iter.Seq[S] {
	// implementation of Heap's algorithm
	// https://en.wikipedia.org/wiki/Heap%27s_algorithm
	return func(yield func(S) bool) {
		k := len(in)
		c := make([]int, k)

		if !yield(in) {
			return
		}

		if k < 2 {
			// no reordering
			return
		}

		i := 1

		for i < k {
			if c[i] < i {
				if i%2 == 0 {
					in[0], in[i] = in[i], in[0]
				} else {
					in[c[i]], in[i] = in[i], in[c[i]]
				}

				if !yield(in) {
					return
				}

				c[i] += 1
				i = 1
			} else {
				c[i] = 0
				i += 1
			}
		}
	}
}
