// package mapper is collection of small mapping helpers.
package mapper

import (
	"fmt"
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/adapter"
)

// Collect maps seq by [slices.Collect].
func Collect[V any](seq iter.Seq[iter.Seq[V]]) iter.Seq[[]V] {
	return adapter.Map(slices.Collect, seq)
}

// Collect2 maps seq by [slices.Collect].
func Collect2[K, V any](seq iter.Seq2[iter.Seq[K], iter.Seq[V]]) iter.Seq2[[]K, []V] {
	return adapter.Map2(
		func(k iter.Seq[K], v iter.Seq[V]) ([]K, []V) {
			return slices.Collect(k), slices.Collect(v)
		},
		seq,
	)
}

// Clone maps seq by [slices.Clone].
func Clone[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[S] {
	return adapter.Map(slices.Clone, seq)
}

// Clone2 maps seq by [slices.Clone].
func Clone2[S1 ~[]E1, S2 ~[]E2, E1, E2 any](seq iter.Seq2[S1, S2]) iter.Seq2[S1, S2] {
	return adapter.Map2(
		func(s1 S1, s2 S2) (S1, S2) { return slices.Clone(s1), slices.Clone(s2) },
		seq,
	)
}

// Sprintf maps seq by [fmt.Sprintf]
func Sprintf[V any](format string, seq iter.Seq[V]) iter.Seq[string] {
	return adapter.Map(
		func(v V) string {
			return fmt.Sprintf(format, v)
		},
		seq,
	)
}

// Sprintf2 maps seq by [fmt.Sprintf]
//
// format receives k as fist and v as second.
// If callers need reversed order, callers should use %[1]verb and %[2]verb.
// see https://pkg.go.dev/fmt@go1.23.3#hdr-Explicit_argument_indexes.
func Sprintf2[K, V any](format string, seq iter.Seq2[K, V]) iter.Seq[string] {
	return hiter.Unify(
		func(k K, v V) string {
			return fmt.Sprintf(format, k, v)
		},
		seq,
	)
}
