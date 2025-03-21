package bufioiter_test

import (
	"bufio"
	"io"
	"iter"
	"strings"
	"testing"
	"testing/iotest"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/bufioiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
	"gotest.tools/v3/assert"
)

var (
	scannerFactory = func() *bufio.Scanner {
		return bufio.NewScanner(strings.NewReader("foo\nbar\nbaz\n"))
	}
	scannerErrFactory = func() *bufio.Scanner {
		return bufio.NewScanner(
			io.MultiReader(
				strings.NewReader("foo\nbar\nbaz\n"),
				iotest.ErrReader(testcase.ErrSample)),
		)
	}
)

func TestScan(t *testing.T) {
	testcase.One[string]{
		Seq: func() iter.Seq[string] {
			return bufioiter.Scan(scannerFactory())
		},
		Seqs: []func() iter.Seq[string]{
			func() iter.Seq[string] {
				return iterable.Scanner{Scanner: scannerFactory()}.IntoIter()
			},
		},
		Expected: []string{"foo", "bar", "baz"},
		BreakAt:  2,
		Stateful: true,
	}.Test(t)

	var scanner *bufio.Scanner
	testcase.One[string]{
		Seq: func() iter.Seq[string] {
			scanner = scannerErrFactory()
			return bufioiter.Scan(scanner)
		},
		Seqs: []func() iter.Seq[string]{
			func() iter.Seq[string] {
				scanner = scannerErrFactory()
				return iterable.Scanner{Scanner: scanner}.IntoIter()
			},
		},
		Expected: []string{"foo", "bar", "baz"},
		BreakAt:  2,
		CmpOpt:   []goCmp.Option{testcase.CompareErrorsIs},
		Stateful: true,
	}.Test(t, func(length, count int) {
		if length != 2 {
			assert.ErrorIs(t, scanner.Err(), testcase.ErrSample)
		}
	})
}

func TestScanErr(t *testing.T) {
	testcase.Two[string, error]{
		Seq: func() iter.Seq2[string, error] {
			return bufioiter.ScanErr(scannerFactory())
		},
		Expected: []hiter.KeyValue[string, error]{
			{K: "foo", V: nil},
			{K: "bar", V: nil},
			{K: "baz", V: nil},
		},
		BreakAt:  2,
		Stateful: true,
	}.Test(t)

	testcase.Two[string, error]{
		Seq: func() iter.Seq2[string, error] {
			return bufioiter.ScanErr(scannerErrFactory())
		},
		Expected: []hiter.KeyValue[string, error]{
			{K: "foo", V: nil},
			{K: "bar", V: nil},
			{K: "baz", V: nil},
			{K: "", V: testcase.ErrSample},
		},
		BreakAt:  2,
		CmpOpt:   []goCmp.Option{testcase.CompareErrorsIs},
		Stateful: true,
	}.Test(t)
}
