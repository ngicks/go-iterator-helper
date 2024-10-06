package sh

import (
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func GroupFunc[U any, K comparable, V any](f func(U) (K, V), seq iter.Seq[U]) map[K][]V {
	return hiter.ReduceGroup(
		func(accum []V, cur V) []V { return append(accum, cur) },
		[]V(nil),
		hiter.Divide(f, seq),
	)
}
