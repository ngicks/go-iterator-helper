package hiter

import (
	"database/sql"
	"iter"
)

// SqlRows returns an iterator over scanned rows from r.
// scanner will be called against [*sql.Rows] after each time [*sql.Rows.Next] returns true.
// It must either call [*sql.Rows.Scan] once per invocation or do nothing and return.
// If the scan result or [*sql.Rows.Err] returns a non-nil error,
// the iterator stops its iteration immediately after yielding the error.
func SqlRows[T any](r *sql.Rows, scanner func(*sql.Rows) (T, error)) iter.Seq2[T, error] {
	return Nexter(r, scanner)
}

// Nexter is like [SqlRows] but extends the input to arbitrary implementors, e.g. sqlx.
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
