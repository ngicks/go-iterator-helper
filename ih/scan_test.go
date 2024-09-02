package ih_test

import (
	"bufio"
	"errors"
	"io"
	"iter"
	"strings"
	"testing"
	"testing/iotest"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/ih"
)

func TestScan(t *testing.T) {
	testCase2[string, error]{
		Seq: func() iter.Seq2[string, error] {
			scanner := bufio.NewScanner(strings.NewReader("foo\nbar\nbaz\n"))
			return ih.Scan(scanner)
		},
		Expected: []ih.KeyValue[string, error]{{"foo", nil}, {"bar", nil}, {"baz", nil}},
		BreakAt:  2,
	}.Test(t)

	sampleErr := errors.New("sample")
	testCase2[string, error]{
		Seq: func() iter.Seq2[string, error] {
			scanner := bufio.NewScanner(
				io.MultiReader(
					strings.NewReader("foo\nbar\nbaz\n"),
					iotest.ErrReader(sampleErr)),
			)
			return ih.Scan(scanner)
		},
		Expected: []ih.KeyValue[string, error]{{"foo", nil}, {"bar", nil}, {"baz", nil}, {"", sampleErr}},
		BreakAt:  2,
		CmpOpt:   []goCmp.Option{goCmp.Comparer(func(e1, e2 error) bool { return errors.Is(e1, e2) })},
	}.Test(t)
}
