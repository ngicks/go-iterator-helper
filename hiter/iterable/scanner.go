package iterable

import (
	"bufio"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var _ hiter.IntoIterable[string] = (*Scanner)(nil)

type Scanner struct {
	*bufio.Scanner
}

func (s Scanner) IntoIter() iter.Seq[string] {
	return hiter.Scan(s.Scanner)
}
