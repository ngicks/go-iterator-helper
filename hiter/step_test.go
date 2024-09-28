package hiter_test

import (
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/x/exp/xiter"
	"gotest.tools/v3/assert"
)

func TestStep(t *testing.T) {
	assert.DeepEqual(
		t,
		[]int{10, 14, 18, 22, 26, 30, 34, 38, 42, 46},
		slices.Collect(xiter.Limit(hiter.Step(10, 4), 10)),
	)
	assert.DeepEqual(
		t,
		[]int{12, 5, -2},
		slices.Collect(xiter.Limit(hiter.Step(12, -7), 3)),
	)
}

func TestStepBy(t *testing.T) {
	src := slices.Collect(
		xiter.Map(func(i int) string {
			return fmt.Sprintf("%d", i)
		},
			hiter.Range(0, 15),
		),
	)
	testCase2[int, string]{
		Seq:      func() iter.Seq2[int, string] { return hiter.StepBy(3, 4, src) },
		Expected: []hiter.KeyValue[int, string]{{3, "3"}, {7, "7"}, {11, "11"}},
		BreakAt:  2,
	}.Test(t)
}