package iterable

import (
	"bytes"
	"cmp"
	"container/list"
	"container/ring"
	"crypto/rand"
	"encoding/hex"
	"io"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func randStr() string {
	buf := new(bytes.Buffer)
	_, _ = io.CopyN(buf, rand.Reader, 8)
	return hex.EncodeToString(buf.Bytes())
}

func TestReverse(t *testing.T) {
	srcL := list.New()
	srcR := ring.New(5)
	srcM := map[int]int{}
	type sample struct {
		Key   int
		Other string
	}
	srcMCmp := map[sample]int{}
	srcS := slices.Collect(hiter.Range(0, 5))

	for i := range 5 {
		// list
		srcL.PushBack(i)
		// ring
		srcR.Value = i
		srcR = srcR.Next()
		// map
		srcM[i] = i
		// map sorted
		srcMCmp[sample{Key: i, Other: randStr()}] = i
	}

	l := ListAll[int]{srcL}
	r := RingAll[int]{srcR}
	m := MapSorted[int, int](srcM)
	mCmp := MapSortedFunc[map[sample]int, sample, int]{
		M:   srcMCmp,
		Cmp: func(i, j sample) int { return cmp.Compare(i.Key, j.Key) },
	}
	s := SliceAll[int](srcS)
	ra := Range[int]{0, 5}

	assertReverse(t, l, l.Reverse())
	assertReverse(t, l.Reverse(), l.Reverse().Reverse())

	assertReverse(t, r, r.Reverse())
	assertReverse(t, r.Reverse(), r.Reverse().Reverse())

	assertReverse2(t, m, m.Reverse())
	assertReverse2(t, m.Reverse(), m.Reverse().Reverse())

	assertReverse2(t, mCmp, mCmp.Reverse())
	assertReverse2(t, mCmp.Reverse(), mCmp.Reverse().Reverse())

	assertReverse(t, s, s.Reverse())
	assertReverse(t, s.Reverse(), s.Reverse().Reverse())

	assertReverse(t, ra, ra.Reverse())
	assertReverse(t, ra.Reverse(), ra.Reverse().Reverse())
}

func assertReverse[V any](t *testing.T, iter hiter.Iterable[V], rev hiter.Iterable[V]) {
	t.Helper()
	result := slices.Collect(iter.Iter())
	slices.Reverse(result)
	reveredResult := slices.Collect(rev.Iter())
	assert.DeepEqual(t, result, reveredResult)
}

func assertReverse2[K, V any](t *testing.T, iter hiter.Iterable2[K, V], rev hiter.Iterable2[K, V]) {
	t.Helper()
	result := hiter.Collect2(iter.Iter2())
	reveredResult := hiter.Collect2(rev.Iter2())
	slices.Reverse(result)
	assert.DeepEqual(t, result, reveredResult)
}
