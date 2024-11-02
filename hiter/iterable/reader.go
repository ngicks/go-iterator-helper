package iterable

import (
	"io"
	"iter"
)

// Reader returns a reader which reads values from seq marshaled by marshaler.
// seq must be stateful and one-shot; each time Read is called one or more values are yielded from seq.
// If seq is pure and reuseable, Read reads same value repeatedly.
// You can wrap pure seq to be stateful by wrapping it with [NewResumable] or [NewPeekable].
func Reader[V any](marshaler func(V) ([]byte, error), seq iter.Seq[V]) io.Reader {
	return &iterReader[V]{
		marshaler: marshaler,
		seq:       seq,
	}
}

type iterReader[V any] struct {
	marshaler func(V) ([]byte, error)
	seq       iter.Seq[V]
	buf       []byte
	err       error
}

func (r *iterReader[V]) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	if len(r.buf) > 0 {
		nn := copy(p, r.buf)
		n += nn
		p = p[nn:]
		r.buf = r.buf[nn:]
		if len(p) == 0 {
			return
		}
	}
	for v := range r.seq {
		r.buf, err = r.marshaler(v)
		if err != nil {
			r.err = err
			return
		}
		nn := copy(p, r.buf)
		n += nn
		p = p[nn:]
		r.buf = r.buf[nn:]
		if len(p) == 0 {
			return
		}
	}
	err = io.EOF
	r.err = err
	return
}
