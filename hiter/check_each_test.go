package hiter_test

import (
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestCheckEach(t *testing.T) {
	src := slices.Collect(hiter.Range(5, 10))

	for i := range hiter.Range(0, 2) {
		var (
			checkReceived  []hiter.KeyValue[int, int]
			checkReceived2 []hiter.KeyValue[hiter.KeyValue[int, int], int]
		)
		result := slices.Collect(
			hiter.CheckEach(
				i,
				func(v, i int) bool {
					checkReceived = append(checkReceived, hiter.KeyValue[int, int]{v, i})
					return true
				},
				slices.Values(src),
			),
		)
		result2 := hiter.Collect2(
			hiter.CheckEach2(
				i,
				func(k, v, i int) bool {
					checkReceived2 = append(checkReceived2, hiter.KeyValue[hiter.KeyValue[int, int], int]{hiter.KeyValue[int, int]{k, v}, i})
					return true
				},
				slices.All(src),
			),
		)
		assert.DeepEqual(t, []int{5, 6, 7, 8, 9}, result)
		assert.DeepEqual(t, []hiter.KeyValue[int, int]{{5, 0}, {6, 1}, {7, 2}, {8, 3}, {9, 4}}, checkReceived)
		assert.DeepEqual(t, []hiter.KeyValue[int, int]{{0, 5}, {1, 6}, {2, 7}, {3, 8}, {4, 9}}, result2)
		assert.DeepEqual(
			t,
			[]hiter.KeyValue[hiter.KeyValue[int, int], int]{
				{hiter.KeyValue[int, int]{0, 5}, 0},
				{hiter.KeyValue[int, int]{1, 6}, 1},
				{hiter.KeyValue[int, int]{2, 7}, 2},
				{hiter.KeyValue[int, int]{3, 8}, 3},
				{hiter.KeyValue[int, int]{4, 9}, 4},
			},
			checkReceived2,
		)
	}

	type testCase struct {
		name  string
		count int
		n     int
	}

	src = slices.Collect(hiter.Range(0, 20))
	src2 := hiter.Collect2(slices.All(src))

	for _, tc := range []testCase{
		{"3*3", 3, 3},
		{"4*3", 4, 3},
		{"2*5", 2, 5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				count          int
				checkReceived  []hiter.KeyValue[int, int]
				checkReceived2 []hiter.KeyValue[hiter.KeyValue[int, int], int]
			)
			count = 0
			result := slices.Collect(
				hiter.CheckEach(
					tc.n,
					func(v, i int) bool {
						checkReceived = append(checkReceived, hiter.KeyValue[int, int]{v, i})
						count++
						return count < tc.count
					},
					slices.Values(src),
				),
			)
			count = 0
			result2 := hiter.Collect2(
				hiter.CheckEach2(
					tc.n,
					func(k, v, i int) bool {
						checkReceived2 = append(checkReceived2, hiter.KeyValue[hiter.KeyValue[int, int], int]{hiter.KeyValue[int, int]{k, v}, i})
						count++
						return count < tc.count
					},
					slices.All(src),
				),
			)

			assert.DeepEqual(t, src[:(tc.n*tc.count)-1], result)
			assert.DeepEqual(t, src2[:(tc.n*tc.count)-1], result2)
		})
	}

}
