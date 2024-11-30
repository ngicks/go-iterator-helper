package mapper

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
)

func TestCancellable(t *testing.T) {
	src := slices.Collect(hiter.Range(0, 20))
	src2 := hiter.Collect2(slices.All(src))

	type testCase struct {
		n           int
		count       int
		shouldYield int
	}

	for _, tc := range []testCase{
		{3, 3, 6},
		{4, 8, 12},
		{2, 5, 6},
	} {
		t.Run(fmt.Sprintf("%#v", tc), func(t *testing.T) {
			t.Run("1", func(t *testing.T) {
				count := tc.count
				ctx, cancel := context.WithCancel(context.Background())
				result := slices.Collect(
					Cancellable(
						tc.n,
						ctx,
						hiter.Tap(
							func(int) {
								if count == 0 {
									cancel()
								}
								count--
							},
							slices.Values(src),
						),
					),
				)
				if len(result) == 0 {
					assert.Assert(t, (tc.shouldYield-1) == 0)
				} else {
					assert.DeepEqual(t, src[:tc.shouldYield-1], result)
				}
			})
		})
		t.Run("2", func(t *testing.T) {
			count := tc.count
			ctx, cancel := context.WithCancel(context.Background())
			result := hiter.Collect2(
				Cancellable2(
					tc.n,
					ctx,
					hiter.Tap2(
						func(int, int) {
							if count == 0 {
								cancel()
							}
							count--
						},
						hiter.Values2(src2),
					),
				),
			)
			if len(result) == 0 {
				assert.Assert(t, (tc.shouldYield-1) == 0)
			} else {
				assert.DeepEqual(t, src2[:tc.shouldYield-1], result)
			}
		})
	}
}

func TestHandleErr(t *testing.T) {
	var handleReceived []int
	result := slices.Collect(
		HandleErr(
			func(i int, err error) bool { handleReceived = append(handleReceived, i); return true },
			hiter.Pairs(hiter.Range(0, 5), hiter.Repeat(error(nil), -1)),
		),
	)
	assert.DeepEqual(t, slices.Collect(hiter.Range(0, 5)), result)
	assert.DeepEqual(t, []int(nil), handleReceived)

	errs := slices.Collect(
		mapIter(
			func(i int64) error { return errors.New(strconv.FormatInt(i, 10)) },
			hiter.Range[int64](1, 6),
		),
	)
	result = slices.Collect(
		HandleErr(
			func(i int, err error) bool {
				handleReceived = append(handleReceived, i)
				return !errors.Is(err, errs[len(errs)-2])
			},
			hiter.Pairs(hiter.Range(0, 6), xiter.Concat(hiter.Once(error(nil)), slices.Values(errs))),
		),
	)
	// values paired to error are excluded.
	assert.DeepEqual(t, []int{0}, result)
	assert.DeepEqual(t, []int{1, 2, 3, 4}, handleReceived)
}
