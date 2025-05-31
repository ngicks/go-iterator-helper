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

var (
	errSample = errors.New("sample")
	expected  = []int{0, 1, 2, 3, 4, 5}
	expected2 = []hiter.KeyValue[int, int]{
		{K: 0, V: 5},
		{K: 1, V: 4},
		{K: 2, V: 3},
		{K: 3, V: 2},
		{K: 4, V: 1},
		{K: 5, V: 0},
	}
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
		Expected: expected,
		BreakAt:  3,
	}.Test(t)

	box := errbox.New(
		hiter.Pairs(
			hiter.Range(0, 6),
			hiter.Concat(hiter.Repeat(error(nil), 5), hiter.Once(errSample)),
		),
	)

	assert.DeepEqual(t, expected[:3], slices.Collect(hiter.Limit(3, box.IntoIter())))
	assert.NilError(t, box.Err())
	assert.DeepEqual(t, expected[:4], slices.Collect(hiter.Limit(4, box.IntoIter())))
	assert.NilError(t, box.Err())
	assert.DeepEqual(t, expected[:5], slices.Collect(box.IntoIter()))
	assert.ErrorIs(t, box.Err(), errSample)
	assert.DeepEqual(t, []int(nil), slices.Collect(box.IntoIter()))
}

func TestBox_Map(t *testing.T) {
	// Test the Map function (not Map2)
	box := errbox.Map(
		func(i int) (int, error) {
			if i >= 3 {
				return 0, errSample
			}
			return i * 2, nil
		},
		hiter.Range(0, 6),
	)
	assert.DeepEqual(t, []int{0, 2, 4}, slices.Collect(box.IntoIter()))
	assert.ErrorIs(t, box.Err(), errSample)
}

func TestBox_Map2(t *testing.T) {
	box := errbox.Map2(
		func(i, j int) (int, int, error) {
			if i > j {
				return j, i, errSample
			}
			return i, j, nil
		},
		hiter.Pairs(hiter.Range(0, 6), hiter.Range(5, -1)),
	)
	assert.DeepEqual(t, expected2[:3], hiter.Collect2(box.IntoIter2()))
	assert.ErrorIs(t, box.Err(), errSample)
}

func TestBox2(t *testing.T) {
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return errbox.New2(
				hiter.Pairs(
					hiter.ToKeyValue(hiter.Pairs(hiter.Range(0, 6), hiter.Range(5, -1))),
					hiter.Repeat(error(nil), -1),
				),
			).IntoIter2()
		},
		Expected: expected2,
		BreakAt:  3,
	}.Test(t)

	box := errbox.New2(
		hiter.Pairs(
			hiter.ToKeyValue(hiter.Pairs(hiter.Range(0, 6), hiter.Range(5, -1))),
			xiter.Concat(hiter.Repeat(error(nil), 5), hiter.Once(errSample)),
		),
	)

	assert.DeepEqual(t, expected2[:3], hiter.Collect2(hiter.Limit2(3, box.IntoIter2())))
	assert.NilError(t, box.Err())
	assert.DeepEqual(t, expected2[:4], hiter.Collect2(hiter.Limit2(4, box.IntoIter2())))
	assert.NilError(t, box.Err())
	assert.DeepEqual(t, expected2[:5], hiter.Collect2(box.IntoIter2()))
	assert.ErrorIs(t, box.Err(), errSample)
	assert.DeepEqual(t, []hiter.KeyValue[int, int](nil), hiter.Collect2(box.IntoIter2()))
}

func TestBox2_Map2(t *testing.T) {
	errSample := errors.New("sample")
	box := errbox.Map2(
		func(i, j int) (int, int, error) {
			if i > j {
				return j, i, errSample
			}
			return i, j, nil
		},
		hiter.Pairs(hiter.Range(0, 6), hiter.Range(5, -1)),
	)
	assert.DeepEqual(t, expected2[:3], hiter.Collect2(box.IntoIter2()))
	assert.ErrorIs(t, box.Err(), errSample)
}

// Mock Nexter implementation for testing
type mockNexter struct {
	data  []int
	index int
	err   error
}

func (m *mockNexter) Next() bool {
	m.index++
	return m.index < len(m.data)
}

func (m *mockNexter) Err() error {
	return m.err
}

func TestNewNexter(t *testing.T) {
	// Test normal case
	mock := &mockNexter{
		data:  []int{1, 2, 3, 4, 5},
		index: -1, // Start at -1 so Next() will increment to 0
		err:   nil,
	}

	nexter := errbox.NewNexter(mock, func(m *mockNexter) (int, error) {
		if m.index >= len(m.data) {
			return 0, errors.New("out of bounds")
		}
		return m.data[m.index], nil
	})

	result := slices.Collect(nexter.IntoIter())
	expected := []int{1, 2, 3, 4, 5}
	assert.DeepEqual(t, expected, result)
	assert.NilError(t, nexter.Err())
}
