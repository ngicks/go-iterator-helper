package hiter_test

import (
	"bufio"
	"errors"
	"io"
	"iter"
	"strings"
	"testing"
	"testing/iotest"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/iterable"
)

func TestScan(t *testing.T) {
	factory := func() *bufio.Scanner {
		return bufio.NewScanner(strings.NewReader("foo\nbar\nbaz\n"))
	}
	testCase2[string, error]{
		Seq: func() iter.Seq2[string, error] {
			return hiter.Scan(factory())
		},
		Seqs: []func() iter.Seq2[string, error]{
			func() iter.Seq2[string, error] {
				return iterable.Scanner{Scanner: factory()}.IntoIter2()
			},
		},
		Expected: []hiter.KeyValue[string, error]{{"foo", nil}, {"bar", nil}, {"baz", nil}},
		BreakAt:  2,
		Stateful: true,
	}.Test(t)

	sampleErr := errors.New("sample")
	factory = func() *bufio.Scanner {
		return bufio.NewScanner(
			io.MultiReader(
				strings.NewReader("foo\nbar\nbaz\n"),
				iotest.ErrReader(sampleErr)),
		)
	}
	testCase2[string, error]{
		Seq: func() iter.Seq2[string, error] {
			return hiter.Scan(factory())
		},
		Seqs: []func() iter.Seq2[string, error]{
			func() iter.Seq2[string, error] {
				return iterable.Scanner{Scanner: factory()}.IntoIter2()
			},
		},
		Expected: []hiter.KeyValue[string, error]{{"foo", nil}, {"bar", nil}, {"baz", nil}, {"", sampleErr}},
		BreakAt:  2,
		CmpOpt:   []goCmp.Option{goCmp.Comparer(func(e1, e2 error) bool { return errors.Is(e1, e2) })},
		Stateful: true,
	}.Test(t)
}
