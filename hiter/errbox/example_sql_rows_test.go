package errbox_test

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ngicks/go-iterator-helper/hiter/errbox"
	"github.com/ngicks/go-iterator-helper/internal/testhelper"
)

func ExampleNewSqlRows_successful() {
	type TestRow struct {
		Id    int
		Title string
		Body  string
	}

	mock := testhelper.OpenMockDB(false)
	defer mock.Close()

	rows, err := mock.Query("SELECT id, title, body FROM posts")
	if err != nil {
		panic(err)
	}

	scanner := func(r *sql.Rows) (TestRow, error) {
		var t TestRow
		err := r.Scan(&t.Id, &t.Title, &t.Body)
		return t, err
	}

	boxed := errbox.NewSqlRows(rows, scanner)

	for row := range boxed.IntoIter() {
		fmt.Printf("row = %#v\n", row)
	}
	fmt.Printf("stored err: %v\n", boxed.Err())
	// Output:
	// row = errbox_test.TestRow{Id:1, Title:"post 1", Body:"hello"}
	// row = errbox_test.TestRow{Id:2, Title:"post 2", Body:"world"}
	// row = errbox_test.TestRow{Id:3, Title:"post 3", Body:"iter"}
	// stored err: <nil>
}

func ExampleNewSqlRows_row_error() {
	mock := testhelper.OpenMockDB(true)
	defer mock.Close()

	boxed := errbox.NewSqlRows(testhelper.QueryRows(mock), testhelper.Scan)

	for row := range boxed.IntoIter() {
		fmt.Printf("row = %#v\n", row)
	}
	fmt.Printf("stored err: %v\n", boxed.Err())
	// Output:
	// row = testhelper.TestRow{Id:1, Title:"post 1", Body:"hello"}
	// row = testhelper.TestRow{Id:2, Title:"post 2", Body:"world"}
	// stored err: mock error
}

func ExampleNewSqlRows_scan_error() {
	scanErr := errors.New("scan")

	mock := testhelper.OpenMockDB(true)
	defer mock.Close()

	var count int
	boxed := errbox.NewSqlRows(
		testhelper.QueryRows(mock),
		func(r *sql.Rows) (testhelper.TestRow, error) {
			count++
			if count > 1 {
				return *new(testhelper.TestRow), scanErr
			}
			return testhelper.Scan(r)
		},
	)

	for row := range boxed.IntoIter() {
		fmt.Printf("row = %#v\n", row)
	}
	fmt.Printf("stored err: %v\n", boxed.Err())
	// Output:
	// row = testhelper.TestRow{Id:1, Title:"post 1", Body:"hello"}
	// stored err: scan
}
