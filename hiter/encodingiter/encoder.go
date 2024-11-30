package encodingiter

import (
	"encoding"
	"io"
	"iter"
)

// Write writes all values from seq by marshaling it through marshaler.
// It returns the number of byte written to w.
// If marshaler or io.Writer returns an error it stops the iteration and the function returns the error with numbers of written bytes so far.
//
// The marshaler receives each value from seq and accumulative number of bytes written.
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

// Write2 is like [Write] but accepts [iter.Seq2].
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

// WriteTextMarshaler writes all calling results of [encoding.TextMarshaler.MarshalText] from seq to w.
// It returns the number of bytes written to w and an error first encountered.
func WriteTextMarshaler[Enc encoding.TextMarshaler](w io.Writer, sep []byte, seq iter.Seq[Enc]) (n int, err error) {
	return Write(
		w,
		func(e Enc, _ int) ([]byte, error) {
			bin, err := e.MarshalText()
			return append(bin, sep...), err
		},
		seq,
	)
}

// WriteBinaryMarshaler writes all calling results of [encoding.BinaryMarshaler.MarshalBinary] from seq to w.
// It returns the number of bytes written to w and an error first encountered.
func WriteBinaryMarshaler[Enc encoding.BinaryMarshaler](w io.Writer, sep []byte, seq iter.Seq[Enc]) (n int, err error) {
	return Write(
		w,
		func(e Enc, _ int) ([]byte, error) {
			bin, err := e.MarshalBinary()
			return append(bin, sep...), err
		},
		seq,
	)
}

// Encode writes consecutive decode results of all values from seq.
//
// The iterator stops if and only if dec returns io.EOF. Handling other errors is caller's responsibility.
// If the first error should stop the iterator, use [LimitUntil], [LimitAfter] or [*errbox.Box].
func Encode[V any, Enc interface{ Encode(v any) error }](enc Enc, seq iter.Seq[V]) error {
	for v := range seq {
		if err := enc.Encode(v); err != nil {
			return err
		}
	}
	return nil
}
