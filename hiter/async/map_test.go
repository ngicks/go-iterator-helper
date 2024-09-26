package async

import (
	"context"
	"errors"
	"iter"
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"time"

	goCmp "github.com/google/go-cmp/cmp"
	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

var (
	compareError = goCmp.Comparer(func(e1, e2 error) bool { return errors.Is(e1, e2) })
)

var (
	errSample = errors.New("sample")
)

func TestMap_successful(t *testing.T) {
	m := func() iter.Seq2[int, error] {
		return Map(
			context.Background(),
			1, 1,
			func(ctx context.Context, i int) (int, error) {
				return i + 10, errSample
			},
			hiter.Range(0, 5),
		)
	}
	result := []hiter.KeyValue[int, error]{
		{K: 10, V: errSample},
		{K: 11, V: errSample},
		{K: 12, V: errSample},
		{K: 13, V: errSample},
		{K: 14, V: errSample},
	}
	assert.Assert(t, cmp.DeepEqual(result, hiter.Collect2(m()), compareError))
	assert.Assert(t, cmp.DeepEqual(result[:2], hiter.Collect2(xiter.Limit2(m(), 2)), compareError))
	assert.Assert(
		t,
		cmp.DeepEqual(
			[]hiter.KeyValue[int, error](nil),
			hiter.Collect2(Map(
				context.Background(),
				1, 1,
				func(ctx context.Context, i int) (int, error) {
					return i + 10, errSample
				},
				hiter.Empty[int](),
			)),
			compareError,
		),
	)
}

func TestMap_limit_blocking_worker(t *testing.T) {
	type testCase struct {
		name        string
		queueLimit  int
		workerLimit int
		actualLimit int // must be less than 5
	}

	for _, tc := range []testCase{
		{"1", 1, 1, 1},
		{"5", 5, 5, 5},
		{"3", 3, 3, 3},
		{"3", 3, 0, 3},
		{"queueLimit<workerLimit", 3, 5, 3},
		{"queueLimit>workerLimit", 2, 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				count         = make(chan struct{}, 100)
				limitConsumer = make(chan struct{})
				resultChan    = make(chan []hiter.KeyValue[int, error], 1)
			)
			go func() {
				resultChan <- hiter.Collect2(Map(
					context.Background(),
					tc.queueLimit, tc.workerLimit,
					func(ctx context.Context, i int) (int, error) {
						count <- struct{}{}
						<-limitConsumer
						// random duration sleep. for random completion order
						time.Sleep(rand.N[time.Duration](100))
						return i + 10, nil
					},
					hiter.Range(0, 5),
				))
			}()

			for range tc.actualLimit {
				<-count
			}
			select {
			case <-time.NewTimer(time.Millisecond).C:
			case <-count:
				t.Fatalf("worker should not be called more than limit")
			}

			close(limitConsumer)

			assert.Assert(
				t,
				cmp.DeepEqual(
					[]hiter.KeyValue[int, error]{
						{K: 10, V: nil},
						{K: 11, V: nil},
						{K: 12, V: nil},
						{K: 13, V: nil},
						{K: 14, V: nil},
					},
					<-resultChan,
					compareError,
				),
			)
		})
	}
}

func TestMap_limit_blocking_consumer(t *testing.T) {
	var (
		count         = make(chan struct{}, 100)
		limitConsumer = make(chan struct{})
		resultChan    = make(chan []hiter.KeyValue[int, error], 1)
	)
	go func() {
		var result []hiter.KeyValue[int, error]
		for k, v := range Map(
			context.Background(),
			3, 1,
			func(ctx context.Context, i int) (int, error) {
				count <- struct{}{}
				return i + 10, nil
			},
			hiter.Range(0, 5),
		) {
			result = append(result, hiter.KeyValue[int, error]{K: k, V: v})
			<-limitConsumer
		}
		resultChan <- result
	}()

	for range 3 {
		<-count
	}
	select {
	case <-time.NewTimer(time.Millisecond).C:
	case <-count:
		t.Fatalf("worker should not be called more than limit")
	}

	for range 2 {
		limitConsumer <- struct{}{}
	}
	for range 2 {
		<-count
	}
	select {
	case <-time.NewTimer(time.Millisecond).C:
	case <-count:
		t.Fatalf("worker should not be called more than limit")
	}

	close(limitConsumer)

	assert.Assert(
		t,
		cmp.DeepEqual(
			[]hiter.KeyValue[int, error]{
				{K: 10, V: nil},
				{K: 11, V: nil},
				{K: 12, V: nil},
				{K: 13, V: nil},
				{K: 14, V: nil},
			},
			<-resultChan,
			compareError,
		),
	)
}

func TestMap_cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var count atomic.Int32
	kv := hiter.Collect2(Map(
		ctx,
		// queueLimit must be exactly 3 to prevent race condition in the test.
		3, 5,
		func(ctx context.Context, i int) (int, error) {
			if count.Add(1) >= 3 {
				cancel()
			}
			<-ctx.Done()
			return i + 10, ctx.Err()
		},
		hiter.Range(0, 5),
	))
	t.Logf("%#v", kv)
	assert.Assert(
		t,
		cmp.DeepEqual(
			[]hiter.KeyValue[int, error]{
				{K: 10, V: context.Canceled},
				{K: 11, V: context.Canceled},
				{K: 12, V: context.Canceled},
			},
			kv,
			compareError,
		),
	)
}

func TestMap_panic_propagation(t *testing.T) {
	t.Run("mapper panics", func(t *testing.T) {
		var count atomic.Int32
		var kv []hiter.KeyValue[int, error]
		func() {
			defer func() {
				rec := recover()
				assert.Assert(t, rec == errSample)
			}()
			for k, v := range Map(
				context.Background(),
				5, 5,
				func(ctx context.Context, i int) (int, error) {
					if count.Add(1) >= 3 {
						panic(errSample)
					}
					<-ctx.Done()
					return i + 10, nil
				},
				hiter.Range(0, 5),
			) {
				kv = append(kv, hiter.KeyValue[int, error]{K: k, V: v})
			}
		}()
		t.Logf("%#v", kv)
		assert.Assert(t, len(kv) == 2)
		for _, kv := range kv {
			assert.NilError(t, kv.V)
		}
	})
	t.Run("seq panics", func(t *testing.T) {
		var count atomic.Int32
		var kv []hiter.KeyValue[int, error]
		func() {
			defer func() {
				rec := recover()
				assert.Assert(t, rec == errSample)
			}()
			for k, v := range Map(
				context.Background(),
				5, 5,
				func(ctx context.Context, i int) (int, error) {
					if count.Add(1) < 3 {
						<-ctx.Done()
					}
					return i + 10, nil
				},
				hiter.Tap(
					func(i int) {
						if i == 3 {
							panic(errSample)
						}
					},
					hiter.Range(0, 5),
				),
			) {
				kv = append(kv, hiter.KeyValue[int, error]{K: k, V: v})
			}
		}()
		t.Logf("%#v", kv)
		assert.Assert(t, len(kv) > 0)
		for _, kv := range kv {
			assert.NilError(t, kv.V)
		}
	})
	t.Run("consumer panics", func(t *testing.T) {
		var kv []hiter.KeyValue[int, error]
		func() {
			defer func() {
				rec := recover()
				assert.Assert(t, rec == errSample)
			}()
			var i int
			for k, v := range Map(
				context.Background(),
				5, 5,
				func(ctx context.Context, i int) (int, error) {
					return i + 10, nil
				},
				hiter.Range(0, 5),
			) {
				kv = append(kv, hiter.KeyValue[int, error]{K: k, V: v})
				if i == 1 {
					panic(errSample)
				}
				i++
			}
		}()
		t.Logf("%#v", kv)
		assert.Assert(t, len(kv) > 0)
		for _, kv := range kv {
			assert.NilError(t, kv.V)
		}
	})
}

func TestMap_param(t *testing.T) {
	for i := -100; i <= 0; i++ {
		func() {
			defer func() {
				assert.Assert(t, recover() != nil)
			}()
			for range Map(
				context.Background(),
				i, 100,
				func(ctx context.Context, i int) (int, error) { return 0, nil },
				hiter.Range(0, 5),
			) {
				//
			}
		}()
	}
	func() {
		defer func() {
			assert.Assert(t, recover() != nil)
		}()
		var ctx context.Context
		for range Map(
			ctx,
			100, 100,
			func(ctx context.Context, i int) (int, error) { return 0, nil },
			hiter.Range(0, 5),
		) {
			//
		}
	}()
	func() {
		defer func() {
			assert.Assert(t, recover() != nil)
		}()
		var mapper func(ctx context.Context, i int) (int, error)
		for range Map(
			context.Background(),
			100, 100,
			mapper,
			hiter.Range(0, 5),
		) {
			//
		}
	}()
}
