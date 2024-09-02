package ih_test

import (
	"iter"
	"testing"

	"github.com/ngicks/go-iterator-helper/ih"
)

func TestKeyValues(t *testing.T) {
	testCase2[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return ih.KeyValues[int, int]{
				{2, 1}, {2, 1}, {0, 4}, {-1, 2},
				{2, 1}, {2, 1}, {1, 4}, {-1, 2},
				{2, 1}, {2, 1}, {2, 6}, {-1, 2},
				{2, 1}, {2, 1}, {3, 9}, {-1, 2},
			}.Iter2()
		},
		Expected: []ih.KeyValue[int, int]{
			{2, 1}, {2, 1}, {0, 4}, {-1, 2},
			{2, 1}, {2, 1}, {1, 4}, {-1, 2},
			{2, 1}, {2, 1}, {2, 6}, {-1, 2},
			{2, 1}, {2, 1}, {3, 9}, {-1, 2},
		},
		BreakAt: 3,
	}.Test(t)
}
