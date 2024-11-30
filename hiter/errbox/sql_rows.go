package errbox

import (
	"database/sql"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/databaseiter"
)

type SqlRows[V any] struct {
	*Box[V]
}

func NewSqlRows[V any](rows *sql.Rows, scanner func(*sql.Rows) (V, error)) *SqlRows[V] {
	return &SqlRows[V]{
		Box: New(databaseiter.SqlRows(rows, scanner)),
	}
}

type Nexter[V any] struct {
	*Box[V]
}

func NewNexter[
	V any,
	N interface {
		Next() bool
		Err() error
	},
](n N, scanner func(N) (V, error)) *Nexter[V] {
	return &Nexter[V]{
		Box: New(hiter.Nexter(n, scanner)),
	}
}
