package iterable

import (
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/ih"
)

var (
	_ ih.Iterable[any]       = SliceAll[any]{}
	_ ih.Iterable2[int, any] = SliceAll[any]{}
)

// SliceAll adds Iter and Iter2 methods to slice of any element E.
// They merely call slices.Values, slices.All respectively.
type SliceAll[E any] []E

func (s SliceAll[E]) Iter() iter.Seq[E] {
	return slices.Values(s)
}

func (s SliceAll[E]) Iter2() iter.Seq2[int, E] {
	return slices.All(s)
}

var (
	_ ih.Iterable2[int, any] = SliceBackward[any]{}
	_ ih.Iterable[any]       = SliceBackward[any]{}
)

// SliceBackward adds Iter and Iter2 methods to slice of any element E.
// They return iterators over []E traversing them backward with descending indices.
type SliceBackward[E any] []E

func (s SliceBackward[E]) Iter() iter.Seq[E] {
	return func(yield func(E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i]) {
				return
			}
		}
	}
}

func (s SliceBackward[E]) Iter2() iter.Seq2[int, E] {
	return slices.Backward(s)
}
