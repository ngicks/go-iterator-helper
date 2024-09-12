package hiter

import (
	"database/sql"
	"iter"
)

// SqlRows returns an iterator over scanned rows from r.
// scanner will be invoked against every rows queried in r.
// scanner should call [*sql.Rows.Scan] once or it can skip the row.
// The returned iterator yields scanned result, including non-nil error.
// If scanner returns an error, or [*sql.Rows.Err] returns non-nil error,
// the iterator yields that error and stops iteration.
func SqlRows[T any](r *sql.Rows, scanner func(*sql.Rows) (T, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for r.Next() {
			t, err := scanner(r)
			if !yield(t, err) {
				return
			}
			if err != nil {
				return
			}
		}
		if r.Err() != nil {
			yield(*new(T), r.Err())
			return
		}
	}
}
