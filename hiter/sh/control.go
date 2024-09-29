package sh

import (
	"context"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Cancellable returns an iterator over seq but it also checks if ctx is cancelled each n elements seq yields.
func Cancellable[V any](n int, ctx context.Context, seq iter.Seq[V]) iter.Seq[V] {
	return hiter.CheckEach(n, func(V, int) bool { return ctx.Err() == nil }, seq)
}

// Cancellable2 returns an iterator over seq but it also checks if ctx is cancelled each n elements seq yields.
func Cancellable2[K, V any](n int, ctx context.Context, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return hiter.CheckEach2(n, func(K, V, int) bool { return ctx.Err() == nil }, seq)
}

// HandleErr returns an iterator over only former value of seq.
// If latter value the seq yields is non-nil then it calls handle.
// If handle returns false the iterator stops.
// Even if handle returns true, values paired to non-nil error are excluded from the returned iterator.
func HandleErr[V any](handle func(V, error) bool, seq iter.Seq2[V, error]) iter.Seq[V] {
	return hiter.OmitL(
		filter2(
			func(_ V, err error) bool { return err == nil },
			hiter.LimitUntil2(func(v V, err error) bool { return err == nil || handle(v, err) }, seq),
		),
	)
}
