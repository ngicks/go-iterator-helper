package testhelper

import (
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
)

var ErrMock = errors.New("mock error")

type TestRow struct {
	Id    int
	Title string
	Body  string
}

func OpenMockDB(lastErr bool) *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world").
		AddRow(3, "post 3", "iter")

	if lastErr {
		rows = rows.RowError(2, ErrMock)
	}

	mock.ExpectQuery("^SELECT (.+) FROM posts$").WillReturnRows(rows)

	return db
}

func QueryRows(mock *sql.DB) *sql.Rows {
	rows, err := mock.Query("SELECT id, title, body FROM posts")
	if err != nil {
		panic(err)
	}
	return rows
}

func Scan(r *sql.Rows) (TestRow, error) {
	var t TestRow
	err := r.Scan(&t.Id, &t.Title, &t.Body)
	return t, err
}
