package iterable

import (
	"bufio"
	"iter"

	"github.com/ngicks/go-iterator-helper/ih"
)

var _ ih.IntoIterable2[string, error] = (*Scanner)(nil)

type Scanner struct {
	*bufio.Scanner
}

func (s Scanner) IntoIter2() iter.Seq2[string, error] {
	return ih.Scan(s.Scanner)
}
