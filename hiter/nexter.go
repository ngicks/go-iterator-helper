package hiter

import (
	"database/sql"
	"iter"
)

var _ = (*sql.Rows)(nil)

// Nexter is tested in databaseiter

// Nexter returns an iterator over Nexter implementation, e.g. [sql.Rows], or even sqlx.Rows.
func Nexter[
	T any,
	Nexter interface {
		Next() bool
		Err() error
	},
](n Nexter, scanner func(Nexter) (T, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for n.Next() {
			t, err := scanner(n)
			if !yield(t, err) {
				return
			}
			if err != nil {
				return
			}
		}
		if n.Err() != nil {
			yield(*new(T), n.Err())
			return
		}
	}
}
