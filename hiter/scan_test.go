package hiter_test

import (
	"bufio"
	"io"
	"iter"
	"strings"
	"testing"
	"testing/iotest"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
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
				iotest.ErrReader(errSample)),
		)
	}
)

func TestScan(t *testing.T) {
	testCase1[string]{
		Seq: func() iter.Seq[string] {
			return hiter.Scan(scannerFactory())
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
	testCase1[string]{
		Seq: func() iter.Seq[string] {
			scanner = scannerErrFactory()
			return hiter.Scan(scanner)
		},
		Seqs: []func() iter.Seq[string]{
			func() iter.Seq[string] {
				scanner = scannerErrFactory()
				return iterable.Scanner{Scanner: scanner}.IntoIter()
			},
		},
		Expected: []string{"foo", "bar", "baz"},
		BreakAt:  2,
		CmpOpt:   []goCmp.Option{compareErrorsIs},
		Stateful: true,
	}.Test(t, func(length, count int) {
		if length != 2 {
			assert.ErrorIs(t, scanner.Err(), errSample)
		}
	})
}

func TestScanErr(t *testing.T) {
	testCase2[string, error]{
		Seq: func() iter.Seq2[string, error] {
			return hiter.ScanErr(scannerFactory())
		},
		Expected: []hiter.KeyValue[string, error]{{"foo", nil}, {"bar", nil}, {"baz", nil}},
		BreakAt:  2,
		Stateful: true,
	}.Test(t)

	testCase2[string, error]{
		Seq: func() iter.Seq2[string, error] {
			return hiter.ScanErr(scannerErrFactory())
		},
		Expected: []hiter.KeyValue[string, error]{
			{"foo", nil},
			{"bar", nil},
			{"baz", nil},
			{"", errSample},
		},
		BreakAt:  2,
		CmpOpt:   []goCmp.Option{compareErrorsIs},
		Stateful: true,
	}.Test(t)
}
