package ih_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"testing"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/ih"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

type testCase1[V any] struct {
	Seq      func() iter.Seq[V]
	Seqs     []func() iter.Seq[V]
	Expected []V
	BreakAt  int
	CmpOpt   []goCmp.Option
}

func (tc testCase1[V]) Test(t *testing.T, cb ...func()) {
	t.Helper()

	for i, seq := range append([](func() iter.Seq[V]){tc.Seq}, tc.Seqs...) {
		t.Run(fmt.Sprintf("#%02d", i), func(t *testing.T) {
			collected := slices.Collect(seq())
			assert.Assert(t, cmp.DeepEqual(collected, tc.Expected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			collected = collected[:tc.BreakAt]

			var i int
			for v := range seq() {
				if i == tc.BreakAt {
					break
				}
				collected[i] = v
				i++
			}
			assert.Assert(t, cmp.DeepEqual(collected, tc.Expected[:tc.BreakAt], tc.CmpOpt...))

			for _, f := range cb {
				f()
			}
		})
	}
}

type testCase2[K, V any] struct {
	Seq      func() iter.Seq2[K, V]
	Seqs     []func() iter.Seq2[K, V]
	Expected []ih.KeyValue[K, V]
	BreakAt  int
	CmpOpt   []goCmp.Option
}

func (tc testCase2[K, V]) Test(t *testing.T, cb ...func()) {
	t.Helper()

	for i, seq := range append([](func() iter.Seq2[K, V]){tc.Seq}, tc.Seqs...) {
		t.Run(fmt.Sprintf("#%02d", i), func(t *testing.T) {

			var collected []ih.KeyValue[K, V]
			for k, v := range seq() {
				collected = append(collected, ih.KeyValue[K, V]{k, v})
			}

			assert.Assert(t, cmp.DeepEqual(collected, tc.Expected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			collected = collected[:tc.BreakAt]
			var i int
			for k, v := range seq() {
				if i == tc.BreakAt {
					break
				}
				collected[i] = ih.KeyValue[K, V]{k, v}
				i++
			}
			assert.Assert(t, cmp.DeepEqual(collected, tc.Expected[:tc.BreakAt], tc.CmpOpt...))

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
}

func (tc testCaseMap[K, V]) Test(t *testing.T, cb ...func()) {
	t.Helper()

	for i, seq := range append([](func() iter.Seq2[K, V]){tc.Seq}, tc.Seqs...) {
		t.Run(fmt.Sprintf("#%02d", i), func(t *testing.T) {

			collected := maps.Collect(seq())

			assert.Assert(t, cmp.DeepEqual(collected, tc.Expected, tc.CmpOpt...))

			for _, f := range cb {
				f()
			}

			var i int
			collected = map[K]V{}
			for k, v := range seq() {
				if i == tc.BreakAt {
					break
				}
				collected[k] = v
				i++
			}
			for k, v := range collected {
				assert.Assert(t, cmp.DeepEqual(v, tc.Expected[k], tc.CmpOpt...), "key=%v", k)
			}
			assert.Assert(t, cmp.Len(collected, tc.BreakAt))

			for _, f := range cb {
				f()
			}
		})
	}
}
