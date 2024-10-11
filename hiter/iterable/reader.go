package iterable

import (
	"io"

	"github.com/ngicks/go-iterator-helper/hiter"
)

func Reader[V any](marshaler func(V) ([]byte, error), seq *Resumable[V]) io.ReadCloser {
	return &iterReader[V]{
		marshaler: marshaler,
		seq:       seq,
	}
}

type iterReader[V any] struct {
	marshaler func(V) ([]byte, error)
	seq       *Resumable[V]
	buf       []byte
}

func (r *iterReader[V]) Read(p []byte) (n int, err error) {
	if len(r.buf) > 0 {
		n = copy(p, r.buf)
		p = p[n:]
		r.buf = r.buf[n:]
		if len(r.buf) > 0 {
			return
		}
	}
	next, ok := hiter.First(r.seq.IntoIter())
	if !ok {
		err = io.EOF
		return
	}
	r.buf, err = r.marshaler(next)
	if err != nil {
		return
	}
	nn := copy(p, r.buf)
	n += nn
	r.buf = r.buf[nn:]
	return
}

func (r *iterReader[V]) Close() error {
	r.seq.Stop()
	return nil
}
