package hiter_test

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
)

func TestForEach(t *testing.T) {
	var num []int
	hiter.ForEach(func(i int) { num = append(num, i) }, hiter.Range(0, 5))
	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 5)), num)
}

func TestForEach2(t *testing.T) {
	var num []hiter.KeyValue[int, int]
	iter := hiter.Pairs(hiter.Range(0, 5), hiter.Range(5, 0))
	hiter.ForEach2(func(i, j int) { num = append(num, hiter.KeyValue[int, int]{i, j}) }, iter)
	assert.DeepEqual(t, hiter.Collect2(iter), num)
}

type fakeErrGroup struct {
	wg      sync.WaitGroup
	err     error
	errOnce sync.Once
}

func (g *fakeErrGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		err := f()
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}
func (g *fakeErrGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

func TestForEachGo(t *testing.T) {
	var (
		g    *fakeErrGroup
		mu   sync.Mutex
		args []int
	)
	g = new(fakeErrGroup)

	err := hiter.ForEachGo(
		context.Background(),
		g,
		func(ctx context.Context, i int) error {
			mu.Lock()
			defer mu.Unlock()
			args = append(args, i)
			return errSample
		},
		hiter.Range(0, 5),
	)
	assert.ErrorIs(t, err, errSample)
	mu.Lock()
	slices.Sort(args)
	mu.Unlock()
	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 5)), args)
}

func TestForEachGo2(t *testing.T) {
	var (
		g    *fakeErrGroup
		mu   sync.Mutex
		args []hiter.KeyValue[int, int]
	)
	g = new(fakeErrGroup)

	err := hiter.ForEachGo2(
		context.Background(),
		g,
		func(ctx context.Context, k, v int) error {
			mu.Lock()
			defer mu.Unlock()
			args = append(args, hiter.KeyValue[int, int]{k, v})
			return errSample
		},
		hiter.Pairs(hiter.Range(0, 5), hiter.Range(5, 0)),
	)
	assert.ErrorIs(t, err, errSample)
	mu.Lock()
	slices.SortFunc(args, func(i, j hiter.KeyValue[int, int]) int { return cmp.Compare(i.K, j.K) })
	mu.Unlock()
	assert.DeepEqual(t, hiter.Collect2(hiter.Pairs(hiter.Range(0, 5), hiter.Range(5, 0))), args)
}

func TestDiscard(t *testing.T) {
	var args []int
	hiter.Discard(hiter.Tap(func(i int) { args = append(args, i) }, hiter.Range(0, 5)))
	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 5)), args)
}

func TestDiscard2(t *testing.T) {
	var args []hiter.KeyValue[int, int]
	src := hiter.Pairs(hiter.Range(0, 5), hiter.Range(5, 0))
	hiter.Discard2(hiter.Tap2(func(i, j int) { args = append(args, hiter.KeyValue[int, int]{i, j}) }, src))
	assert.DeepEqual(t, hiter.Collect2(src), args)
}

var (
	trySrc = hiter.Pairs(
		slices.Values([]string{"foo", "bar", "baz", "qux"}),
		slices.Values([]error{nil, nil, errSample, nil}),
	)
)

func TestTryFind(t *testing.T) {
	var (
		v   string
		idx int
		err error
	)
	v, idx, err = hiter.TryFind(func(s string) bool { return s == "bar" }, trySrc)
	assert.Equal(t, "bar", v)
	assert.Equal(t, 1, idx)
	assert.NilError(t, err)

	v, idx, err = hiter.TryFind(func(s string) bool { return s == "baz" }, trySrc)
	assert.Equal(t, "", v)
	assert.Equal(t, -1, idx)
	assert.ErrorIs(t, err, errSample)

	v, idx, err = hiter.TryFind(func(s string) bool { return s == "baz" }, xiter.Limit2(trySrc, 2))
	assert.Equal(t, "", v)
	assert.Equal(t, -1, idx)
	assert.NilError(t, err)
}

func TestTryForEach(t *testing.T) {
	var (
		args []string
		err  error
	)

	args = args[:0]
	err = hiter.TryForEach(func(s string) { args = append(args, s) }, xiter.Limit2(trySrc, 2))
	assert.DeepEqual(t, slices.Collect(hiter.OmitL(xiter.Limit2(trySrc, 2))), args)
	assert.NilError(t, err)

	args = args[:0]
	err = hiter.TryForEach(func(s string) { args = append(args, s) }, trySrc)
	assert.DeepEqual(t, slices.Collect(hiter.OmitL(xiter.Limit2(trySrc, 2))), args)
	assert.ErrorIs(t, err, errSample)
}

func TestTryReduce(t *testing.T) {
	var (
		sum []string
		err error
	)
	sum, err = hiter.TryReduce(func(ss []string, s string) []string { return append(ss, s) }, []string{}, trySrc)
	assert.DeepEqual(t, slices.Collect(hiter.OmitL(xiter.Limit2(trySrc, 2))), sum)
	assert.ErrorIs(t, err, errSample)

	sum, err = hiter.TryReduce(func(ss []string, s string) []string { return append(ss, s) }, []string{}, xiter.Limit2(trySrc, 1))
	assert.DeepEqual(t, []string{"foo"}, sum)
	assert.NilError(t, err)
}