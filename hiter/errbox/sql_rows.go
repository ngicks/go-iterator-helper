package errbox

import (
	"database/sql"

	"github.com/ngicks/go-iterator-helper/hiter"
)

type SqlRows[V any] struct {
	*Box[V]
}

func NewSqlRows[V any](rows *sql.Rows, scanner func(*sql.Rows) (V, error)) *SqlRows[V] {
	return &SqlRows[V]{
		Box: New(hiter.SqlRows(rows, scanner)),
	}
}
