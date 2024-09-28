package hiter_test

import (
	"database/sql"
	"errors"
	"iter"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
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
			testCase2[testhelper.TestRow, error]{
				Seq: func() iter.Seq2[testhelper.TestRow, error] {
					mock = testhelper.OpenMockDB(false)
					return hiter.SqlRows(testhelper.QueryRows(mock), scanner)
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
					{testhelper.TestRow{Id: 1, Title: "post 1", Body: "hello"}, nil},
					{testhelper.TestRow{Id: 2, Title: "post 2", Body: "world"}, nil},
					{testhelper.TestRow{Id: 3, Title: "post 3", Body: "iter"}, nil},
				},
				BreakAt:  2,
				Stateful: true,
			}.Test(t, func(_, _ int) { _ = mock.Close() })
		})

		t.Run("errbox.SqlRows", func(t *testing.T) {
			var boxed *errbox.SqlRows[testhelper.TestRow]
			defer func() { boxed.Stop() }()

			testCase1[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					if boxed != nil {
						boxed.Stop()
					}

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
			defer func() { boxed.Stop() }()

			testCase1[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					if boxed != nil {
						boxed.Stop()
					}

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
			testCase2[testhelper.TestRow, error]{
				Seq: func() iter.Seq2[testhelper.TestRow, error] {
					mock = testhelper.OpenMockDB(true)
					return hiter.SqlRows(testhelper.QueryRows(mock), scanner)
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
					{testhelper.TestRow{Id: 1, Title: "post 1", Body: "hello"}, nil},
					{testhelper.TestRow{Id: 2, Title: "post 2", Body: "world"}, nil},
					{testhelper.TestRow{}, testhelper.ErrMock},
				},
				BreakAt:  2,
				CmpOpt:   []goCmp.Option{compareError},
				Stateful: true,
			}.Test(t, func(_, _ int) { _ = mock.Close() })
		})

		t.Run("*errbox.SqlRows", func(t *testing.T) {
			var boxed *errbox.SqlRows[testhelper.TestRow]
			defer func() { boxed.Stop() }()
			testCase1[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					if boxed != nil {
						boxed.Stop()
					}

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
			defer func() { boxed.Stop() }()
			testCase1[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					if boxed != nil {
						boxed.Stop()
					}

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
			mock    *sql.DB
			count   int
			mockErr = errors.New("sample")
		)
		scanner := func(r *sql.Rows) (testhelper.TestRow, error) {
			var t testhelper.TestRow
			count++
			if count > 2 {
				return t, mockErr
			}
			err := r.Scan(&t.Id, &t.Title, &t.Body)
			return t, err
		}
		t.Run("hiter.SqlRows", func(t *testing.T) {
			testCase2[testhelper.TestRow, error]{
				Seq: func() iter.Seq2[testhelper.TestRow, error] {
					count = 0
					mock = testhelper.OpenMockDB(false)
					return hiter.SqlRows(testhelper.QueryRows(mock), scanner)
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
					{testhelper.TestRow{Id: 1, Title: "post 1", Body: "hello"}, nil},
					{testhelper.TestRow{Id: 2, Title: "post 2", Body: "world"}, nil},
					{testhelper.TestRow{}, mockErr},
				},
				BreakAt:  1,
				CmpOpt:   []goCmp.Option{compareError},
				Stateful: true,
			}.Test(t, func(_, _ int) { _ = mock.Close() })
		})
		t.Run("*errbox.SqlRows", func(t *testing.T) {
			var boxed *errbox.SqlRows[testhelper.TestRow]
			defer func() { boxed.Stop() }()

			testCase1[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					if boxed != nil {
						boxed.Stop()
					}

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
					assert.ErrorIs(t, boxed.Err(), mockErr)
				}
			})
		})
		t.Run("*errbox.Nexter", func(t *testing.T) {
			var boxed *errbox.Nexter[testhelper.TestRow]
			defer func() { boxed.Stop() }()

			testCase1[testhelper.TestRow]{
				Seq: func() iter.Seq[testhelper.TestRow] {
					if boxed != nil {
						boxed.Stop()
					}

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
					assert.ErrorIs(t, boxed.Err(), mockErr)
				}
			})
		})
	})
}
