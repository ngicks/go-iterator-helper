package iterable

import (
	"database/sql"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

var (
	_ hiter.IntoIterable2[any, error] = SqlRows[any]{}
)

// SqlRows adds IntoIter2 to [*sql.Rows].
// For detailed behavior, see [hiter.SqlRows].
type SqlRows[T any] struct {
	Rows    *sql.Rows
	Scanner func(*sql.Rows) (T, error)
}

func (s SqlRows[T]) IntoIter2() iter.Seq2[T, error] {
	return hiter.SqlRows[T](s.Rows, s.Scanner)
}
