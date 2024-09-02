package ih_test

import (
	"iter"
	"sync"
	"testing"

	"github.com/ngicks/go-iterator-helper/ih"
	"github.com/ngicks/go-iterator-helper/ih/iterable"
)

func TestSyncMap(t *testing.T) {
	expected := map[string]string{
		"foo": "foofoo",
		"bar": "barbar",
		"baz": "bazbaz",
	}

	var m sync.Map
	for k, v := range expected {
		m.Store(k, v)
	}

	testCaseMap[string, string]{
		Seq: func() iter.Seq2[string, string] {
			return ih.SyncMap[string, string](&m)
		},
		Seqs: []func() iter.Seq2[string, string]{
			func() iter.Seq2[string, string] {
				return iterable.SyncMap[string, string]{Map: &m}.Iter2()
			},
		},
		Expected: expected,
		BreakAt:  2,
	}.Test(t)

	testCaseMap[string, string]{
		Seq: func() iter.Seq2[string, string] {
			return ih.SyncMap[string, string](&sync.Map{})
		},
		Seqs: []func() iter.Seq2[string, string]{
			func() iter.Seq2[string, string] {
				return iterable.SyncMap[string, string]{Map: &sync.Map{}}.Iter2()
			},
		},
		Expected: map[string]string{},
	}.Test(t)
}
