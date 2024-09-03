package iterable

import (
	"bufio"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ hiter.IntoIterable2[string, error] = (*Scanner)(nil)

type Scanner struct {
	*bufio.Scanner
}

func (s Scanner) IntoIter2() iter.Seq2[string, error] {
	return hiter.Scan(s.Scanner)
}
