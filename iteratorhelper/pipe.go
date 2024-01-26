package iteratorhelper

import (
	"io"
	"iter"
)

type iterReader struct {
	buf []byte
	next func() ([]byte, error, bool)
	stop func()
}

func (r *iterReader) Read(p []byte) (n int, err error) {
	if len(r.buf) == 0 {
		var ok bool
		r.buf, err, ok = r.next()
		if !ok {
			return 0, io.EOF
		}
		if err != nil {
			r.stop()
			return 0, err
		}
	}
	n = copy(p, r.buf)
	if n > 0 {
		r.buf = r.buf[n:]
	}
	return n, nil
}

func (r *iterReader) Close() error {
	r.stop()
	return nil
}

type yieldWriter func(b []byte, err error) bool

func (w yieldWriter) Write(p []byte) (n int, err error) {
	if !w(p, nil) {
		return 0, io.ErrUnexpectedEOF
	}
	return len(p), nil
}

func (w yieldWriter) CloseWithError(err error) error {
	w(nil, err)
	return nil
}

type CloserWithError interface {
	io.Writer
	CloseWithError(err error) error	
}

func Pipe(fn func(w CloserWithError)) io.ReadCloser {
	next, stop := iter.Pull2(func(yield func(b []byte, err error) bool) {
		fn(yieldWriter(yield))
	})
	return &iterReader{nil, next, stop}
}