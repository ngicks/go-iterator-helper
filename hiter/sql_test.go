package hiter_test

import (
	"database/sql"
	"errors"
	"iter"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestSqliteRows(t *testing.T) {
	mockErr := errors.New("mock error")
	openMockSql := func(lastErr bool) *sql.DB {
		db, mock, err := sqlmock.New()
		if err != nil {
			panic(err)
		}

		rows := sqlmock.NewRows([]string{"id", "title", "body"}).
			AddRow(1, "post 1", "hello").
			AddRow(2, "post 2", "world").
			AddRow(3, "post 3", "iter")

		if lastErr {
			rows = rows.RowError(2, mockErr)
		}

		mock.ExpectQuery("^SELECT (.+) FROM posts$").WillReturnRows(rows)

		return db
	}
	queryRows := func(mock *sql.DB) *sql.Rows {
		rows, err := mock.Query("SELECT id, title, body FROM posts")
		if err != nil {
			panic(err)
		}
		return rows
	}
	type testRow struct {
		Id    int
		Title string
		Body  string
	}

	t.Run("successful", func(t *testing.T) {
		var mock *sql.DB
		scanner := func(r *sql.Rows) (testRow, error) {
			var t testRow
			err := r.Scan(&t.Id, &t.Title, &t.Body)
			return t, err
		}
		testCase2[testRow, error]{
			Seq: func() iter.Seq2[testRow, error] {
				mock = openMockSql(false)
				return hiter.SqlRows(queryRows(mock), scanner)
			},
			Seqs: []func() iter.Seq2[testRow, error]{
				func() iter.Seq2[testRow, error] {
					mock = openMockSql(false)
					return iterable.SqlRows[testRow]{Rows: queryRows(mock), Scanner: scanner}.IntoIter2()
				},
			},
			Expected: []hiter.KeyValue[testRow, error]{
				{testRow{Id: 1, Title: "post 1", Body: "hello"}, nil},
				{testRow{Id: 2, Title: "post 2", Body: "world"}, nil},
				{testRow{Id: 3, Title: "post 3", Body: "iter"}, nil},
			},
			BreakAt:  2,
			Stateful: true,
		}.Test(t, func() { _ = mock.Close() })
	})

	t.Run("row error", func(t *testing.T) {
		var mock *sql.DB
		scanner := func(r *sql.Rows) (testRow, error) {
			var t testRow
			err := r.Scan(&t.Id, &t.Title, &t.Body)
			return t, err
		}
		testCase2[testRow, error]{
			Seq: func() iter.Seq2[testRow, error] {
				mock = openMockSql(true)
				return hiter.SqlRows(queryRows(mock), scanner)
			},
			Seqs: []func() iter.Seq2[testRow, error]{
				func() iter.Seq2[testRow, error] {
					mock = openMockSql(true)
					return iterable.SqlRows[testRow]{Rows: queryRows(mock), Scanner: scanner}.IntoIter2()
				},
			},
			Expected: []hiter.KeyValue[testRow, error]{
				{testRow{Id: 1, Title: "post 1", Body: "hello"}, nil},
				{testRow{Id: 2, Title: "post 2", Body: "world"}, nil},
				{testRow{}, mockErr},
			},
			BreakAt:  2,
			CmpOpt:   []goCmp.Option{compareError},
			Stateful: true,
		}.Test(t, func() { _ = mock.Close() })
	})

	t.Run("scan error", func(t *testing.T) {
		var (
			mock    *sql.DB
			count   int
			mockErr = errors.New("sample")
		)
		scanner := func(r *sql.Rows) (testRow, error) {
			var t testRow
			count++
			if count > 1 {
				return t, mockErr
			}
			err := r.Scan(&t.Id, &t.Title, &t.Body)
			return t, err
		}
		testCase2[testRow, error]{
			Seq: func() iter.Seq2[testRow, error] {
				count = 0
				mock = openMockSql(false)
				return hiter.SqlRows(queryRows(mock), scanner)
			},
			Seqs: []func() iter.Seq2[testRow, error]{
				func() iter.Seq2[testRow, error] {
					count = 0
					mock = openMockSql(false)
					return iterable.SqlRows[testRow]{Rows: queryRows(mock), Scanner: scanner}.IntoIter2()
				},
			},
			Expected: []hiter.KeyValue[testRow, error]{
				{testRow{Id: 1, Title: "post 1", Body: "hello"}, nil},
				{testRow{}, mockErr},
			},
			BreakAt:  1,
			CmpOpt:   []goCmp.Option{compareError},
			Stateful: true,
		}.Test(t, func() { _ = mock.Close() })
	})
}
