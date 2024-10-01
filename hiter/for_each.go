package hiter

import (
	"context"
	"iter"
	"slices"

	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
)

var (
	_ = xiter.Reduce[any, any]
	_ = slices.AppendSeq[[]any, any]
)

// ForEach iterates over seq and calls fn with every value seq yields.
func ForEach[V any](fn func(V), seq iter.Seq[V]) {
	for v := range seq {
		fn(v)
	}
}

// ForEach2 iterates over seq and calls fn with every key-value pair seq yields.
func ForEach2[K, V any](fn func(K, V), seq iter.Seq2[K, V]) {
	for k, v := range seq {
		fn(k, v)
	}
}

type GoGroup interface {
	Go(f func() error)
	Wait() error
}

// ForEachGo iterates over seq and calls fn with every value from seq in G's Go method.
// After all values are consumed, the result of Wait is returned.
// You may want to use [*errgroup.Group](https://pkg.go.dev/golang.org/x/sync/errgroup#Group) as implementor.
func ForEachGo[V any, G GoGroup](ctx context.Context, g G, fn func(context.Context, V) error, seq iter.Seq[V]) error {
	for v := range seq {
		g.Go(func() error {
			return fn(ctx, v)
		})
	}
	return g.Wait()
}

// ForEachGo2 iterates over seq and calls fn with every key-value pair from seq in G's Go method.
// After all values are consumed, the result of Wait is returned.
// You may want to use [*errgroup.Group](https://pkg.go.dev/golang.org/x/sync/errgroup#Group) as implementor.
func ForEachGo2[K, V any, G GoGroup](ctx context.Context, g G, fn func(context.Context, K, V) error, seq iter.Seq2[K, V]) error {
	for k, v := range seq {
		g.Go(func() error {
			return fn(ctx, k, v)
		})
	}
	return g.Wait()
}

// Discard fully consumes seq without doing anything.
func Discard[V any](seq iter.Seq[V]) {
	for range seq {
	}
}

// Discard2 fully consumes seq without doing anything.
func Discard2[K, V any](seq iter.Seq2[K, V]) {
	for range seq {
	}
}

// TryFind is like [FindFunc] but stops if seq yields non-nil error.
func TryFind[V any](f func(V) bool, seq iter.Seq2[V, error]) (v V, idx int, err error) {
	var i int
	for v, err := range seq {
		if err != nil {
			return *new(V), -1, err
		}
		if f(v) {
			return v, i, nil
		}
		i++
	}
	return *new(V), -1, nil
}

// TryForEach is like [ForEach] but returns an error if seq yields non-nil error.
func TryForEach[V any](f func(V), seq iter.Seq2[V, error]) error {
	for v, err := range seq {
		if err != nil {
			return err
		}
		f(v)
	}
	return nil
}

// TryReduce is like [xiter.Reduce] but returns an error if seq yields non-nil error.
func TryReduce[Sum, V any](f func(Sum, V) Sum, sum Sum, seq iter.Seq2[V, error]) (Sum, error) {
	for v, err := range seq {
		if err != nil {
			return sum, err
		}
		sum = f(sum, v)
	}
	return sum, nil
}

// TryCollect is like [slices.Collect] but stops collecting at the first error and
// returns the extended result before the error.
func TryCollect[E any](seq iter.Seq2[E, error]) ([]E, error) {
	return TryAppendSeq[[]E](nil, seq)
}

// TryAppendSeq is like [slices.AppendSeq] but stops collecting at the first error
// and returns the extended result before the error.
func TryAppendSeq[S ~[]E, E any](s S, seq iter.Seq2[E, error]) (S, error) {
	for e, err := range seq {
		if err != nil {
			return s, err
		}
		s = append(s, e)
	}
	return s, nil
}
