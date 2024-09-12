package hiter_test

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

var (
	compareError = goCmp.Comparer(func(e1, e2 error) bool { return errors.Is(e1, e2) })
)

type testCase1[V any] struct {
	Seq      func() iter.Seq[V]
	Seqs     []func() iter.Seq[V]
	Expected []V
	BreakAt  int
	CmpOpt   []goCmp.Option
	Stateful bool
}

func (tc testCase1[V]) Test(t *testing.T, cb ...func()) {
	t.Helper()

	for i, seq := range append([](func() iter.Seq[V]){tc.Seq}, tc.Seqs...) {
		t.Run(fmt.Sprintf("#%02d", i), func(t *testing.T) {
			t.Helper()
			collected := slices.Collect(seq())
			assert.Assert(t, cmp.DeepEqual(tc.Expected, collected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			collected = collected[:tc.BreakAt]

			var i int
			s := seq()
			for v := range s {
				if i == tc.BreakAt {
					break
				}
				collected[i] = v
				i++
			}
			assert.Assert(t, cmp.DeepEqual(tc.Expected[:tc.BreakAt], collected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			// call seq after breaking it to check if it is idempotence.
			collected2 := slices.Collect(s)
			if !tc.Stateful {
				assert.Assert(t, cmp.DeepEqual(tc.Expected, collected2))
			} else {
				assert.Assert(t, len(tc.Expected) != len(collected2), "expected not to be=\n%#v\n\nactual=\n%#v\n", tc.Expected, collected2)
			}

			for _, f := range cb {
				f()
			}
		})
	}
}

type testCase2[K, V any] struct {
	Seq      func() iter.Seq2[K, V]
	Seqs     []func() iter.Seq2[K, V]
	Expected []hiter.KeyValue[K, V]
	BreakAt  int
	CmpOpt   []goCmp.Option
	Stateful bool
}

func (tc testCase2[K, V]) Test(t *testing.T, cb ...func()) {
	t.Helper()

	for i, seq := range append([](func() iter.Seq2[K, V]){tc.Seq}, tc.Seqs...) {
		t.Run(fmt.Sprintf("#%02d", i), func(t *testing.T) {
			t.Helper()
			var collected []hiter.KeyValue[K, V]
			for k, v := range seq() {
				collected = append(collected, hiter.KeyValue[K, V]{k, v})
			}

			assert.Assert(t, cmp.DeepEqual(tc.Expected, collected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			collected = collected[:tc.BreakAt]
			var i int
			s := seq()
			for k, v := range s {
				if i == tc.BreakAt {
					break
				}
				collected[i] = hiter.KeyValue[K, V]{k, v}
				i++
			}
			assert.Assert(t, cmp.DeepEqual(tc.Expected[:tc.BreakAt], collected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			// call seq after breaking it to check if it is idempotence.
			collected2 := hiter.Collect2(s)
			if !tc.Stateful {
				assert.Assert(t, cmp.DeepEqual(tc.Expected, collected2))
			} else {
				assert.Assert(t, len(tc.Expected) != len(collected2), "expected not to be=\n%#v\n\nactual=\n%#v\n", tc.Expected, collected2)
			}

			for _, f := range cb {
				f()
			}
		})
	}
}

type testCaseMap[K comparable, V any] struct {
	Seq      func() iter.Seq2[K, V]
	Seqs     []func() iter.Seq2[K, V]
	Expected map[K]V
	BreakAt  int
	CmpOpt   []goCmp.Option
	Stateful bool
}

func (tc testCaseMap[K, V]) Test(t *testing.T, cb ...func()) {
	t.Helper()

	for i, seq := range append([](func() iter.Seq2[K, V]){tc.Seq}, tc.Seqs...) {
		t.Run(fmt.Sprintf("#%02d", i), func(t *testing.T) {
			t.Helper()

			collected := maps.Collect(seq())

			assert.Assert(t, cmp.DeepEqual(tc.Expected, collected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			var i int
			collected = map[K]V{}
			s := seq()
			for k, v := range s {
				if i == tc.BreakAt {
					break
				}
				collected[k] = v
				i++
			}
			for k, v := range collected {
				assert.Assert(t, cmp.DeepEqual(tc.Expected[k], v, tc.CmpOpt...), "key=%v", k)
			}
			assert.Assert(t, cmp.Len(collected, tc.BreakAt))

			for _, f := range cb {
				f()
			}

			// call seq after breaking it to check if it is idempotence.
			collected2 := maps.Collect(s)
			if !tc.Stateful {
				assert.Assert(t, cmp.DeepEqual(tc.Expected, collected2))
			} else {
				assert.Assert(t, len(tc.Expected) != len(collected2), "expected not to be=\n%#v\n\nactual=\n%#v\n", tc.Expected, collected2)
			}

			for _, f := range cb {
				f()
			}
		})
	}
}
