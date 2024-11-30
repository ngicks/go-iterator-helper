package databaseiter

import (
	"database/sql"
	"iter"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// SqlRows returns an iterator over scanned rows from r.
// scanner will be called against [*sql.SqlRows] after each time [*sql.SqlRows.Next] returns true.
// It must either call [*sql.SqlRows.Scan] once per invocation or do nothing and return.
// If the scan result or [*sql.SqlRows.Err] returns a non-nil error,
// the iterator stops its iteration immediately after yielding the error.
func SqlRows[T any](r *sql.Rows, scanner func(*sql.Rows) (T, error)) iter.Seq2[T, error] {
	return hiter.Nexter(r, scanner)
}
