package databaseiter_test

import (
	"database/sql"
	"iter"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/databaseiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"github.com/ngicks/go-iterator-helper/internal/testhelper"
	"gotest.tools/v3/assert"
)

func TestSqliteRows(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		var mock *sql.DB
		scanner := func(r *sql.Rows) (testhelper.TestRow, error) {
			var t testhelper.TestRow
			err := r.Scan(&t.Id, &t.Title, &t.Body)
			return t, err
		}
		t.Run("hiter.SqlRows", func(t *testing.T) {
			testcase.Two[testhelper.TestRow, error]{
				Seq: func() iter.Seq2[testhelper.TestRow, error] {
					mock = testhelper.OpenMockDB(false)
					return databaseiter.SqlRows(testhelper.QueryRows(mock), scanner)
				},
				Seqs: []func() iter.Seq2[testhelper.TestRow, error]{
					func() iter.Seq2[testhelper.TestRow, error] {
						mock = testhelper.OpenMockDB(false)
						return iterable.SqlRows[testhelper.TestRow]{Rows: testhelper.QueryRows(mock), Scanner: scanner}.IntoIter2()
					},
					func() iter.Seq2[testhelper.TestRow, error] {
						mock = testhelper.OpenMockDB(false)
						return hiter.Nexter(testhelper.QueryRows(mock), scanner)
					},
					func() iter.Seq2[testhelper.TestRow, error] {
						mock = testhelper.OpenMockDB(false)
						return iterable.Nexter[testhelper.TestRow, *sql.Rows]{Nexter: testhelper.QueryRows(mock), Scanner: scanner}.IntoIter2()
					},
				},
				Expected: []hiter.KeyValue[testhelper.TestRow, error]{
					{K: testhelper.TestRow{Id: 1, Title: "post 1", Body: "hello"}, V: nil},
					{K: testhelper.TestRow{Id: 2, Title: "post 2", Body: "world"}, V: nil},
					{K: testhelper.TestRow{Id: 3, Title: "post 3", Body: "iter"}, V: nil},
				},
				BreakAt:  2,
				Stateful: true,
			}.Test(t, func(_, _ int) { _ = mock.Close() })
		})

		t.Run("errbox.SqlRows", func(t *testing.T) {
			var boxed *errbox.SqlRows[testhelper.TestRow]

			testcase.One[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					mock = testhelper.OpenMockDB(false)
					boxed = errbox.NewSqlRows(testhelper.QueryRows(mock), scanner)
					return boxed.IntoIter()
				},
				Expected: []testhelper.TestRow{
					{Id: 1, Title: "post 1", Body: "hello"},
					{Id: 2, Title: "post 2", Body: "world"},
					{Id: 3, Title: "post 3", Body: "iter"},
				},
				BreakAt:  1,
				Stateful: true,
			}.Test(t, func(_, _ int) {
				assert.NilError(t, boxed.Err())
			})
		})
		t.Run("*errbox.Nexter", func(t *testing.T) {
			var boxed *errbox.Nexter[testhelper.TestRow]

			testcase.One[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					mock = testhelper.OpenMockDB(false)
					boxed = errbox.NewNexter(testhelper.QueryRows(mock), scanner)
					return boxed.IntoIter()
				},
				Expected: []testhelper.TestRow{
					{Id: 1, Title: "post 1", Body: "hello"},
					{Id: 2, Title: "post 2", Body: "world"},
					{Id: 3, Title: "post 3", Body: "iter"},
				},
				BreakAt:  1,
				Stateful: true,
			}.Test(t, func(_, _ int) {
				assert.NilError(t, boxed.Err())
			})
		})
	})

	t.Run("row error", func(t *testing.T) {
		var mock *sql.DB
		scanner := func(r *sql.Rows) (testhelper.TestRow, error) {
			var t testhelper.TestRow
			err := r.Scan(&t.Id, &t.Title, &t.Body)
			return t, err
		}
		t.Run("hiter.SqlRows", func(t *testing.T) {
			testcase.Two[testhelper.TestRow, error]{
				Seq: func() iter.Seq2[testhelper.TestRow, error] {
					mock = testhelper.OpenMockDB(true)
					return databaseiter.SqlRows(testhelper.QueryRows(mock), scanner)
				},
				Seqs: []func() iter.Seq2[testhelper.TestRow, error]{
					func() iter.Seq2[testhelper.TestRow, error] {
						mock = testhelper.OpenMockDB(true)
						return iterable.SqlRows[testhelper.TestRow]{Rows: testhelper.QueryRows(mock), Scanner: scanner}.IntoIter2()
					},
					func() iter.Seq2[testhelper.TestRow, error] {
						mock = testhelper.OpenMockDB(true)
						return hiter.Nexter(testhelper.QueryRows(mock), scanner)
					},
					func() iter.Seq2[testhelper.TestRow, error] {
						mock = testhelper.OpenMockDB(true)
						return iterable.Nexter[testhelper.TestRow, *sql.Rows]{Nexter: testhelper.QueryRows(mock), Scanner: scanner}.IntoIter2()
					},
				},
				Expected: []hiter.KeyValue[testhelper.TestRow, error]{
					{K: testhelper.TestRow{Id: 1, Title: "post 1", Body: "hello"}, V: nil},
					{K: testhelper.TestRow{Id: 2, Title: "post 2", Body: "world"}, V: nil},
					{K: testhelper.TestRow{}, V: testhelper.ErrMock},
				},
				BreakAt:  2,
				CmpOpt:   []goCmp.Option{testcase.CompareErrorsIs},
				Stateful: true,
			}.Test(t, func(_, _ int) { _ = mock.Close() })
		})

		t.Run("*errbox.SqlRows", func(t *testing.T) {
			var boxed *errbox.SqlRows[testhelper.TestRow]
			testcase.One[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					mock = testhelper.OpenMockDB(true)
					boxed = errbox.NewSqlRows(testhelper.QueryRows(mock), scanner)
					return boxed.IntoIter()
				},
				Expected: []testhelper.TestRow{
					{Id: 1, Title: "post 1", Body: "hello"},
					{Id: 2, Title: "post 2", Body: "world"},
				},
				BreakAt:  1,
				Stateful: true,
			}.Test(t, func(_, count int) {
				if count == 1 {
					assert.NilError(t, boxed.Err())
				} else {
					assert.ErrorIs(t, boxed.Err(), testhelper.ErrMock)
				}
			})
		})
		t.Run("*errbox.Nexter", func(t *testing.T) {
			var boxed *errbox.Nexter[testhelper.TestRow]
			testcase.One[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					mock = testhelper.OpenMockDB(true)
					boxed = errbox.NewNexter(testhelper.QueryRows(mock), scanner)
					return boxed.IntoIter()
				},
				Expected: []testhelper.TestRow{
					{Id: 1, Title: "post 1", Body: "hello"},
					{Id: 2, Title: "post 2", Body: "world"},
				},
				BreakAt:  1,
				Stateful: true,
			}.Test(t, func(_, count int) {
				if count == 1 {
					assert.NilError(t, boxed.Err())
				} else {
					assert.ErrorIs(t, boxed.Err(), testhelper.ErrMock)
				}
			})
		})
	})

	t.Run("scan error", func(t *testing.T) {
		var (
			mock  *sql.DB
			count int
		)
		scanner := func(r *sql.Rows) (testhelper.TestRow, error) {
			var t testhelper.TestRow
			count++
			if count > 2 {
				return t, testcase.ErrSample
			}
			err := r.Scan(&t.Id, &t.Title, &t.Body)
			return t, err
		}
		t.Run("hiter.SqlRows", func(t *testing.T) {
			testcase.Two[testhelper.TestRow, error]{
				Seq: func() iter.Seq2[testhelper.TestRow, error] {
					count = 0
					mock = testhelper.OpenMockDB(false)
					return databaseiter.SqlRows(testhelper.QueryRows(mock), scanner)
				},
				Seqs: []func() iter.Seq2[testhelper.TestRow, error]{
					func() iter.Seq2[testhelper.TestRow, error] {
						count = 0
						mock = testhelper.OpenMockDB(false)
						return iterable.SqlRows[testhelper.TestRow]{Rows: testhelper.QueryRows(mock), Scanner: scanner}.IntoIter2()
					},
					func() iter.Seq2[testhelper.TestRow, error] {
						count = 0
						mock = testhelper.OpenMockDB(false)
						return hiter.Nexter(testhelper.QueryRows(mock), scanner)
					},
					func() iter.Seq2[testhelper.TestRow, error] {
						count = 0
						mock = testhelper.OpenMockDB(false)
						return iterable.Nexter[testhelper.TestRow, *sql.Rows]{Nexter: testhelper.QueryRows(mock), Scanner: scanner}.IntoIter2()
					},
				},
				Expected: []hiter.KeyValue[testhelper.TestRow, error]{
					{K: testhelper.TestRow{Id: 1, Title: "post 1", Body: "hello"}, V: nil},
					{K: testhelper.TestRow{Id: 2, Title: "post 2", Body: "world"}, V: nil},
					{K: testhelper.TestRow{}, V: testcase.ErrSample},
				},
				BreakAt:  1,
				CmpOpt:   []goCmp.Option{testcase.CompareErrorsIs},
				Stateful: true,
			}.Test(t, func(_, _ int) { _ = mock.Close() })
		})
		t.Run("*errbox.SqlRows", func(t *testing.T) {
			var boxed *errbox.SqlRows[testhelper.TestRow]

			testcase.One[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					count = 0
					mock = testhelper.OpenMockDB(false)
					boxed = errbox.NewSqlRows(testhelper.QueryRows(mock), scanner)
					return boxed.IntoIter()
				},
				Expected: []testhelper.TestRow{
					{Id: 1, Title: "post 1", Body: "hello"},
					{Id: 2, Title: "post 2", Body: "world"},
				},
				BreakAt:  1,
				Stateful: true,
			}.Test(t, func(_, count int) {
				if count == 1 {
					assert.NilError(t, boxed.Err())
				} else {
					assert.ErrorIs(t, boxed.Err(), testcase.ErrSample)
				}
			})
		})
		t.Run("*errbox.Nexter", func(t *testing.T) {
			var boxed *errbox.Nexter[testhelper.TestRow]

			testcase.One[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					count = 0
					mock = testhelper.OpenMockDB(false)
					boxed = errbox.NewNexter(testhelper.QueryRows(mock), scanner)
					return boxed.IntoIter()
				},
				Expected: []testhelper.TestRow{
					{Id: 1, Title: "post 1", Body: "hello"},
					{Id: 2, Title: "post 2", Body: "world"},
				},
				BreakAt:  1,
				Stateful: true,
			}.Test(t, func(_, count int) {
				if count == 1 {
					assert.NilError(t, boxed.Err())
				} else {
					assert.ErrorIs(t, boxed.Err(), testcase.ErrSample)
				}
			})
		})
	})
}
