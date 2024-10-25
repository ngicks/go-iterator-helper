package iterable

import (
	"io"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Reader returns a reader which reads values from seq marshaled by marshaler.
// seq must be one-shot; each time Read is called one or two values are yielded from seq.
// If seq is pure and reuseable, Read reads same value repeatedly.
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
		n = copy(p, r.buf)
		p = p[n:]
		r.buf = r.buf[n:]
		if len(r.buf) > 0 {
			return
		}
	}
	next, ok := hiter.First(r.seq)
	if !ok {
		err = io.EOF
		r.err = err
		return
	}
	r.buf, err = r.marshaler(next)
	if err != nil {
		r.err = err
		return
	}
	nn := copy(p, r.buf)
	n += nn
	r.buf = r.buf[nn:]
	return
}
