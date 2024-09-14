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
	"gotest.tools/v3/assert"
)

func TestScan(t *testing.T) {
	factory := func() *bufio.Scanner {
		return bufio.NewScanner(strings.NewReader("foo\nbar\nbaz\n"))
	}
	testCase1[string]{
		Seq: func() iter.Seq[string] {
			return hiter.Scan(factory())
		},
		Seqs: []func() iter.Seq[string]{
			func() iter.Seq[string] {
				return iterable.Scanner{Scanner: factory()}.IntoIter()
			},
		},
		Expected: []string{"foo", "bar", "baz"},
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
	var scanner *bufio.Scanner
	testCase1[string]{
		Seq: func() iter.Seq[string] {
			scanner = factory()
			return hiter.Scan(scanner)
		},
		Seqs: []func() iter.Seq[string]{
			func() iter.Seq[string] {
				scanner = factory()
				return iterable.Scanner{Scanner: scanner}.IntoIter()
			},
		},
		Expected: []string{"foo", "bar", "baz"},
		BreakAt:  2,
		CmpOpt:   []goCmp.Option{compareError},
		Stateful: true,
	}.Test(t, func(length, count int) {
		if length != 2 {
			assert.ErrorIs(t, scanner.Err(), sampleErr)
		}
	})
}
