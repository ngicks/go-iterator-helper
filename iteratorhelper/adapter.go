package iteratorhelper

import (
	"iter"
)

func Chain[K, V any](
	iters ...(func(yield func(k K, v V) bool)),
) func(yield func(k K, v V) bool) {
	return func(yield func(k K, v V) bool) {
		stopped := false
		for _, iter := range iters {
			iter(func(k K, v V) bool {
				if !yield(k, v) {
					stopped = true
					return false
				}
				return true
			})
			if stopped {
				return
			}
		}
	}
}

func Chunk[K, V any](iter func(yield func(k K, v V) bool), size uint) func(yield func(k []K, v []V) bool) {
	if size == 0 {
		panic("Chunk: size must not be zero.")
	}
	return func(yield func(k []K, v []V) bool) {
		stopped := false
		chunkedK, chunkedV := make([]K, size), make([]V, size)
		idx := uint(0)
		iter(func(k K, v V) bool {
			if idx < size {
				chunkedK[idx] = k
				chunkedV[idx] = v
				idx++
			}
			if idx >= size {
				idx = 0
				if !yield(append([]K{}, chunkedK...), append([]V{}, chunkedV...)) {
					stopped = true
					return false
				}
			}
			return true
		})
		if !stopped && idx != 0 {
			yield(append([]K{}, chunkedK[:idx]...), append([]V{}, chunkedV[:idx]...))
		}
	}
}

func Enumerate[V any](iter func(yield func(v V) bool)) func(yield func(idx int, v V) bool) {
	return func(yield func(idx int, v V) bool) {
		var idx int
		iter(func(v V) bool {
			if !yield(idx, v) {
				return false
			}
			idx++
			return true
		})
	}
}

func FilterSelect[K, V any](
	iter func(yield func(k K, v V) bool),
	selector func(k K, v V) bool,
) func(yield func(k K, v V) bool) {
	return func(yield func(k K, v V) bool) {
		iter(func(k K, v V) bool {
			if selector(k, v) && !yield(k, v) {
				return false
			}
			return true
		})
	}
}

func FilterExclude[K, V any](
	iter func(yield func(k K, v V) bool),
	excluder func(k K, v V) bool,
) func(yield func(k K, v V) bool) {
	return func(yield func(k K, v V) bool) {
		iter(func(k K, v V) bool {
			if !excluder(k, v) && !yield(k, v) {
				return false
			}
			return true
		})
	}
}

func Map[K1, V1, K2, V2 any](
	iter func(yield func(k K1, v V1) bool),
	mapper func(k K1, v V1) (K2, V2),
) func(yield func(k K2, v V2) bool) {
	return func(yield func(k K2, v V2) bool) {
		iter(func(k K1, v V1) bool {
			if !yield(mapper(k, v)) {
				return false
			}
			return true
		})
	}
}

func SkipWhile[K, V any](
	iter func(yield func(k K, v V) bool),
	predicate func(k K, v V) bool,
) func(yield func(k K, v V) bool) {
	return func(yield func(k K, v V) bool) {
		taking := false
		iter(func(k K, v V) bool {
			if !taking && !predicate(k, v) {
				taking = true
			}
			if taking && !yield(k, v) {
				return false
			}
			return true
		})
	}
}

func TakeWhile[K, V any](
	iter func(yield func(k K, v V) bool),
	predicate func(k K, v V) bool,
) func(yield func(k K, v V) bool) {
	return func(yield func(k K, v V) bool) {
		iter(func(k K, v V) bool {
			if !predicate(k, v) {
				return false
			}
			if !yield(k, v) {
				return false
			}
			return true
		})
	}
}

func Window[K, V any](iter func(yield func(k K, v V) bool), size uint) func(yield func(k []K, v []V) bool) {
	if size == 0 {
		panic("Window: size must not be zero.")
	}
	return func(yield func(k []K, v []V) bool) {
		bufK, bufV := make([]K, size), make([]V, size)
		idx := uint(0)
		ended := false
		iter(func(k K, v V) bool {
			if idx < size {
				bufK[idx] = k
				bufV[idx] = v
				idx++
			} else {
				copy(bufK, bufK[1:])
				copy(bufV, bufV[1:])
				bufK[len(bufK)-1] = k
				bufV[len(bufV)-1] = v
			}

			if idx == size {
				if !yield(append([]K{}, bufK...), append([]V{}, bufV...)) {
					ended = true
					return false
				}
			}

			return true
		})

		if !ended && idx != size {
			_ = yield(append([]K{}, bufK[:idx]...), append([]V{}, bufV[:idx]...))
		}
	}
}

func Zip[V1, V2 any](
	left func(yield func(v V1) bool),
	right func(yield func(v V2) bool),
) func(yield func(l V1, r V2) bool) {
	return func(yield func(l V1, r V2) bool) {
		iterResult := make(chan bool)
		leftCh := make(chan V1)
		rightCh := make(chan V2)

		go func() {
			left(func(v V1) bool {
				leftCh <- v
				return <-iterResult
			})
			close(leftCh)
		}()
		go func() {
			right(func(v V2) bool {
				rightCh <- v
				return <-iterResult
			})
			close(rightCh)
		}()

		defer func() {
			close(iterResult)
			for _ = range leftCh {
			}
			for _ = range rightCh {
			}
		}()

		for {
			l, okL := <-leftCh
			r, okR := <-rightCh
			if !okL || !okR {
				break
			}
			if !yield(l, r) {
				break
			}
			iterResult <- true
			iterResult <- true
		}
	}
}

func ZipPull[V1, V2 any](
	left func(yield func(v V1) bool),
	right func(yield func(v V2) bool),
) func(yield func(l V1, r V2) bool) {
	return func(yield func(l V1, r V2) bool) {
		nextL, stopL := iter.Pull(left)
		nextR, stopR := iter.Pull(right)
		defer stopL()
		defer stopR()

		for {
			l, lOk := nextL()
			r, rOk := nextR()

			if !lOk || !rOk {
				return
			}
			if !yield(l, r) {
				return
			}
		}
	}
}

type Pair[K, V any] struct {
	K K
	V V
}

func ZipPair[K1, V1, K2, V2 any](
	left func(yield func(k K1, v V1) bool),
	right func(yield func(k K2, v V2) bool),
) func(yield func(l Pair[K1, V1], r Pair[K2, V2]) bool) {
	return func(yield func(l Pair[K1, V1], r Pair[K2, V2]) bool) {
		iterResult := make(chan bool)
		leftCh := make(chan Pair[K1, V1])
		rightCh := make(chan Pair[K2, V2])

		go func() {
			left(func(k K1, v V1) bool {
				leftCh <- Pair[K1, V1]{K: k, V: v}
				return <-iterResult
			})
			close(leftCh)
		}()
		go func() {
			right(func(k K2, v V2) bool {
				rightCh <- Pair[K2, V2]{K: k, V: v}
				return <-iterResult
			})
			close(rightCh)
		}()

		defer func() {
			close(iterResult)
			for _ = range leftCh {
			}
			for _ = range rightCh {
			}
		}()

		for {
			l, okL := <-leftCh
			r, okR := <-rightCh
			if !okL || !okR {
				break
			}
			if !yield(l, r) {
				break
			}
			iterResult <- true
			iterResult <- true
		}
	}
}

func Swap[K, V any](iter func(yield func(k K, v V) bool)) func(yield func(v V, k K) bool) {
	return func(yield func(v V, k K) bool) {
		iter(func(k K, v V) bool {
			if !yield(v, k) {
				return false
			}
			return true
		})
	}
}
