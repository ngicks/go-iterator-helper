package tee

import (
	_ "io"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

// TeeSeq is [iter.Seq] equivalent of [io.TeeReader].
//
// TeeSeq returns a [iter.Seq] that pushes to pusher what it reads from seq.
// Yielding values from the returned iterator performs push before the inner loop receives the value.
// The iterator is not stateful; you may want to wrap it with [iterable.Resumable].
// If pusher returns false, the iterator stops iteration without yielding value.
//
// Experimental: not tested and might be changed any time.
func TeeSeq[V any](pusher func(v V) bool, seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !pusher(v) {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// TeeSeq2 is [iter.Seq2] equivalent of [io.TeeReader].
//
// TeeSeq2 returns a [iter.Seq2] that pushes to pusher what it reads from seq.
// Yielding key-value pairs from the returned iterator performs push before the inner loop receives the pair.
// The iterator is not stateful; you may want to wrap it with [iterable.Resumable2].
// If pusher returns false, the iterator stops iteration without yielding pair.
func TeeSeq2[K, V any](pusher func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if !pusher(k, v) {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// TeeSeqPipe tees values from seq to [*Pipe]. see doc comments for [TeeSeq].
// Yielding values from returned [*iterable.Resumable] also performs push to [*Pipe].
func TeeSeqPipe[V any](bufSize int, seq iter.Seq[V]) (*Pipe[V], *iterable.Resumable[V]) {
	p := NewPipe[V](bufSize)
	tee := iterable.NewResumable(
		hiter.TapLast(
			p.Close,
			TeeSeq(p.Push, seq),
		),
	)
	return p, tee
}

// TeeSeqPipe2 tees key-value pairs from seq to [*Pipe2]. see doc comments for [TeeSeq2].
// Yielding pairs from returned [*iterable.Resumable2] also performs push to [*Pipe2].
func TeeSeqPipe2[K, V any](bufSize int, seq iter.Seq2[K, V]) (*Pipe2[K, V], *iterable.Resumable2[K, V]) {
	p := NewPipe2[K, V](bufSize)
	tee := iterable.NewResumable2(
		hiter.TapLast2(
			p.Close,
			TeeSeq2(p.Push, seq),
		),
	)
	return p, tee
}
