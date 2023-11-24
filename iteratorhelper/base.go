package iteratorhelper

import (
	"bufio"
	"io"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func SliceIter[T ~[]E, E any](sl T) func(yield func(k int, v E) bool) {
	return func(yield func(k int, v E) bool) {
		for idx, v := range sl {
			if !yield(idx, v) {
				return
			}
		}
	}
}

func SliceIterSingle[T ~[]E, E any](sl T) func(yield func(v E) bool) {
	return func(yield func(v E) bool) {
		for _, v := range sl {
			if !yield(v) {
				return
			}
		}
	}
}

func MapIter[K comparable, V any](m map[K]V) func(yield func(k K, v V) bool) {
	return func(yield func(k K, v V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func ChanIter[V any](ch <-chan V) func(yield func(V) bool) {
	return func(yield func(V) bool) {
		for v := range ch {
			if !yield(v) {
				return
			}
		}
	}
}

func RangeIter(start, end int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func OrderedMapIter[K comparable, V any](o *orderedmap.OrderedMap[K, V]) func(yield func(k K, v V) bool) {
	return func(yield func(k K, v V) bool) {
		for pair := o.Oldest(); pair != nil; pair = pair.Next() {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

func Scan(r io.Reader, split bufio.SplitFunc) func(yield func(text string, err error) bool) {
	scanner := bufio.NewScanner(r)
	if split != nil {
		scanner.Split(split)
	}

	return func(yield func(text string, err error) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text(), nil) {
				return
			}
		}
		if scanner.Err() != nil {
			yield("", scanner.Err())
			return
		}
	}
}
