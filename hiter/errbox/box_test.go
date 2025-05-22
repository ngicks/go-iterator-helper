package errbox_test

import (
	"errors"
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/errbox"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
)

func TestBox(t *testing.T) {
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return errbox.New(
				hiter.Pairs(
					hiter.Range(0, 6),
					hiter.Repeat(error(nil), -1),
				),
			).IntoIter()
		},
		Expected: []int{0, 1, 2, 3, 4, 5},
		BreakAt:  3,
	}.Test(t)

	errSample := errors.New("sample")
	box := errbox.New(
		hiter.Pairs(
			hiter.Range(0, 6),
			xiter.Concat(hiter.Repeat(error(nil), 5), hiter.Once(errSample)),
		),
	)

	assert.DeepEqual(t, []int{0, 1, 2}, slices.Collect(hiter.Limit(3, box.IntoIter())))
	assert.NilError(t, box.Err())
	assert.DeepEqual(t, []int{0, 1, 2, 3}, slices.Collect(hiter.Limit(4, box.IntoIter())))
	assert.NilError(t, box.Err())
	assert.DeepEqual(t, []int{0, 1, 2, 3, 4}, slices.Collect(box.IntoIter()))
	assert.ErrorIs(t, box.Err(), errSample)
	assert.DeepEqual(t, []int(nil), slices.Collect(box.IntoIter()))
}
