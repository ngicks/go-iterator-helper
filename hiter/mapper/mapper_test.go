package mapper

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestCollect(t *testing.T) {
	assert.DeepEqual(
		t,
		[][]int{{0, 1, 2}, {1, 2, 3}, {2, 3, 4}},
		slices.Collect(Collect(hiter.WindowSeq(3, hiter.Range(0, 5)))),
	)
}

func TestCollect2(t *testing.T) {
	assert.DeepEqual(
		t,
		[]hiter.KeyValue[[]int, []int]{
			{K: []int{0, 1, 2}, V: []int{5, 4, 3}},
			{K: []int{1, 2, 3}, V: []int{4, 3, 2}},
			{K: []int{2, 3, 4}, V: []int{3, 2, 1}},
		},
		hiter.Collect2(
			Collect2(
				hiter.Pairs(
					hiter.WindowSeq(3, hiter.Range(0, 5)),
					hiter.WindowSeq(3, hiter.Range(5, 0)),
				),
			),
		),
	)
}

func TestClone(t *testing.T) {
	src := slices.Collect(hiter.Range(0, 5))

	result := slices.Collect(Clone(hiter.Window(src, 3)))
	assert.DeepEqual(t, [][]int{{0, 1, 2}, {1, 2, 3}, {2, 3, 4}}, result)
	result[0][1] = 50
	assert.DeepEqual(t, [][]int{{0, 50, 2}, {1, 2, 3}, {2, 3, 4}}, result)

	// Double check the effect doing same thing to result without Clone.
	resultWoClone := slices.Collect(hiter.Window(src, 3))
	assert.DeepEqual(t, [][]int{{0, 1, 2}, {1, 2, 3}, {2, 3, 4}}, resultWoClone)
	resultWoClone[0][1] = 50
	assert.DeepEqual(t, [][]int{{0, 50, 2}, {50, 2, 3}, {2, 3, 4}}, resultWoClone)
}

func TestClone2(t *testing.T) {
	src1 := slices.Collect(hiter.Range(0, 5))
	src2 := slices.Collect(hiter.Range(5, 0))

	result := hiter.Collect2(Clone2(hiter.Pairs(hiter.Window(src1, 3), hiter.Window(src2, 3))))
	assert.DeepEqual(t,
		[]hiter.KeyValue[[]int, []int]{
			{K: []int{0, 1, 2}, V: []int{5, 4, 3}},
			{K: []int{1, 2, 3}, V: []int{4, 3, 2}},
			{K: []int{2, 3, 4}, V: []int{3, 2, 1}},
		},
		result,
	)
	result[0].K[1] = 50
	assert.DeepEqual(t,
		[]hiter.KeyValue[[]int, []int]{
			{K: []int{0, 50, 2}, V: []int{5, 4, 3}},
			{K: []int{1, 2, 3}, V: []int{4, 3, 2}},
			{K: []int{2, 3, 4}, V: []int{3, 2, 1}},
		},
		result,
	)

	// Double check the effect doing same thing to result without Clone.
	resultWoClone := hiter.Collect2(hiter.Pairs(hiter.Window(src1, 3), hiter.Window(src2, 3)))
	assert.DeepEqual(t,
		[]hiter.KeyValue[[]int, []int]{
			{K: []int{0, 1, 2}, V: []int{5, 4, 3}},
			{K: []int{1, 2, 3}, V: []int{4, 3, 2}},
			{K: []int{2, 3, 4}, V: []int{3, 2, 1}},
		},
		resultWoClone,
	)
	resultWoClone[0].K[1] = 50
	assert.DeepEqual(t,
		[]hiter.KeyValue[[]int, []int]{
			{K: []int{0, 50, 2}, V: []int{5, 4, 3}},
			{K: []int{50, 2, 3}, V: []int{4, 3, 2}},
			{K: []int{2, 3, 4}, V: []int{3, 2, 1}},
		},
		resultWoClone,
	)
}

func TestSprintf(t *testing.T) {
	assert.DeepEqual(
		t,
		[]string{"000", "001", "002"},
		slices.Collect(Sprintf("%03d", hiter.Range(0, 3))),
	)
}

func TestSprintf2(t *testing.T) {
	assert.DeepEqual(
		t,
		[]string{"002, 000", "001, 001", "000, 002"},
		slices.Collect(Sprintf2("%03[2]d, %03[1]d", hiter.Enumerate(hiter.Range(2, -1)))),
	)
}

func TestStringsReplacer(t *testing.T) {
	replacer := NewStringsReplacer("old", "new", "foo", "bar")

	input := []string{"old text", "foo test", "something"}
	result := slices.Collect(replacer.Map(slices.Values(input)))
	expected := []string{"new text", "bar test", "something"}

	assert.DeepEqual(t, expected, result)
}

func TestReplacer(t *testing.T) {
	replacer := Replacer[string, int]{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	input := []string{"one", "two", "three", "four"}
	result := slices.Collect(replacer.Map(slices.Values(input)))
	expected := []int{1, 2, 3, 0} // "four" maps to zero value of int

	assert.DeepEqual(t, expected, result)
}
