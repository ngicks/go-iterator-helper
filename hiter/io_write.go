package hiter

import (
	"io"
	"iter"
)

func Write[V any](w io.Writer, marshaler func(v V, written int) ([]byte, error), seq iter.Seq[V]) (n int, er error) {
	for v := range seq {
		bin, err := marshaler(v, n)
		if err != nil {
			return n, err
		}
		nn, err := w.Write(bin)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func Write2[K, V any](w io.Writer, marshaler func(k K, v V, written int) ([]byte, error), seq iter.Seq2[K, V]) (n int, er error) {
	for k, v := range seq {
		bin, err := marshaler(k, v, n)
		if err != nil {
			return n, err
		}
		nn, err := w.Write(bin)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func Encode[V any, Enc interface{ Encode(v any) error }](enc Enc, seq iter.Seq[V]) error {
	for v := range seq {
		if err := enc.Encode(v); err != nil {
			return err
		}
	}
	return nil
}
