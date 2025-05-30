package hiter_test

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
	"gotest.tools/v3/assert"
	goCmp "gotest.tools/v3/assert/cmp"
)

func TestConcat(t *testing.T) {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := slices.Values([]int{4, 5, 6})
	seq3 := slices.Values([]int{7, 8})

	t.Run("multiple sequences", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Concat(seq1, seq2, seq3)
			},
			Expected: []int{1, 2, 3, 4, 5, 6, 7, 8},
			BreakAt:  4,
		}.Test(t)
	})

	t.Run("single sequence", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Concat(seq1)
			},
			Expected: []int{1, 2, 3},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("empty sequences", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Concat[int]()
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})

	t.Run("with empty sequence in middle", func(t *testing.T) {
		emptySeq := slices.Values([]int{})
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Concat(seq1, emptySeq, seq2)
			},
			Expected: []int{1, 2, 3, 4, 5, 6},
			BreakAt:  3,
		}.Test(t)
	})
}

func TestConcat2(t *testing.T) {
	seq1 := slices.All([]string{"a", "b"})
	seq2 := slices.All([]string{"c", "d"})
	seq3 := slices.All([]string{"e"})

	t.Run("multiple sequences", func(t *testing.T) {
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Concat2(seq1, seq2, seq3)
			},
			Expected: []hiter.KeyValue[int, string]{
				{0, "a"},
				{1, "b"},
				{0, "c"},
				{1, "d"},
				{0, "e"},
			},
			BreakAt: 3,
		}.Test(t)
	})

	t.Run("empty sequences", func(t *testing.T) {
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Concat2[int, string]()
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		seq1     []int
		seq2     []int
		expected bool
	}{
		{"identical sequences", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"different sequences", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"different lengths", []int{1, 2}, []int{1, 2, 3}, false},
		{"empty sequences", []int{}, []int{}, true},
		{"one empty", []int{1}, []int{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hiter.Equal(slices.Values(tt.seq1), slices.Values(tt.seq2))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEqual2(t *testing.T) {
	tests := []struct {
		name     string
		seq1     []hiter.KeyValue[string, int]
		seq2     []hiter.KeyValue[string, int]
		expected bool
	}{
		{
			"identical sequences",
			[]hiter.KeyValue[string, int]{{"a", 1}, {"b", 2}},
			[]hiter.KeyValue[string, int]{{"a", 1}, {"b", 2}},
			true,
		},
		{
			"different values",
			[]hiter.KeyValue[string, int]{{"a", 1}, {"b", 2}},
			[]hiter.KeyValue[string, int]{{"a", 1}, {"b", 3}},
			false,
		},
		{
			"different length",
			[]hiter.KeyValue[string, int]{{"a", 1}, {"b", 2}},
			[]hiter.KeyValue[string, int]{{"a", 1}},
			false,
		},
		{
			"different keys",
			[]hiter.KeyValue[string, int]{{"a", 1}, {"b", 2}},
			[]hiter.KeyValue[string, int]{{"a", 1}, {"c", 2}},
			false,
		},
		{
			"empty sequences",
			[]hiter.KeyValue[string, int]{},
			[]hiter.KeyValue[string, int]{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hiter.Equal2(hiter.Values2(tt.seq1), hiter.Values2(tt.seq2))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEqualFunc(t *testing.T) {
	tests := []struct {
		name     string
		seq1     []string
		seq2     []int
		cmpFunc  func(string, int) bool
		expected bool
	}{
		{
			"equal by length",
			[]string{"hi", "bye"},
			[]int{10, 100},
			func(s string, i int) bool { return len(s) == len(fmt.Sprint(i)) },
			true,
		},
		{
			"not equal by length",
			[]string{"hi", "hello"},
			[]int{10, 100},
			func(s string, i int) bool { return len(s) == len(fmt.Sprint(i)) },
			false,
		},
		{
			"empty sequences",
			[]string{},
			[]int{},
			func(s string, i int) bool { return true },
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hiter.EqualFunc(slices.Values(tt.seq1), slices.Values(tt.seq2), tt.cmpFunc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEqualFunc2(t *testing.T) {
	tests := []struct {
		name     string
		seq1     []hiter.KeyValue[string, int]
		seq2     []hiter.KeyValue[int, string]
		cmpFunc  func(string, int, int, string) bool
		expected bool
	}{
		{
			"equal by swapped key-value",
			[]hiter.KeyValue[string, int]{{"1", 1}, {"2", 2}},
			[]hiter.KeyValue[int, string]{{1, "1"}, {2, "2"}},
			func(k1 string, v1 int, k2 int, v2 string) bool {
				return k1 == v2 && v1 == k2
			},
			true,
		},
		{
			"not equal",
			[]hiter.KeyValue[string, int]{{"1", 1}, {"2", 3}},
			[]hiter.KeyValue[int, string]{{1, "1"}, {2, "2"}},
			func(k1 string, v1 int, k2 int, v2 string) bool {
				return k1 == v2 && v1 == k2
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hiter.EqualFunc2(hiter.Values2(tt.seq1), hiter.Values2(tt.seq2), tt.cmpFunc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilter(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	t.Run("even numbers", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Filter(func(v int) bool { return v%2 == 0 }, slices.Values(src))
			},
			Expected: []int{2, 4, 6, 8, 10},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("greater than 5", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Filter(func(v int) bool { return v > 5 }, slices.Values(src))
			},
			Expected: []int{6, 7, 8, 9, 10},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("no matches", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Filter(func(v int) bool { return v > 100 }, slices.Values(src))
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})

	t.Run("all matches", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Filter(func(v int) bool { return v > 0 }, slices.Values(src))
			},
			Expected: src,
			BreakAt:  5,
		}.Test(t)
	})
}

func TestFilter2(t *testing.T) {
	src := []string{"hello", "world", "foo", "bar", "baz"}

	t.Run("filter by index and value", func(t *testing.T) {
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Filter2(
					func(k int, v string) bool { return k%2 == 0 && len(v) > 3 },
					slices.All(src),
				)
			},
			Expected: []hiter.KeyValue[int, string]{{0, "hello"}},
			BreakAt:  1,
		}.Test(t)
	})

	t.Run("filter by value length", func(t *testing.T) {
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Filter2(
					func(k int, v string) bool { return len(v) == 3 },
					slices.All(src),
				)
			},
			Expected: []hiter.KeyValue[int, string]{{2, "foo"}, {3, "bar"}, {4, "baz"}},
			BreakAt:  2,
		}.Test(t)
	})
}

func TestLimit(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	t.Run("limit 5", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Limit(5, slices.Values(src))
			},
			Expected: []int{1, 2, 3, 4, 5},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("limit 0", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Limit(0, slices.Values(src))
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})

	t.Run("limit negative", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Limit(-1, slices.Values(src))
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})

	t.Run("limit greater than length", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Limit(20, slices.Values(src))
			},
			Expected: src,
			BreakAt:  5,
		}.Test(t)
	})
}

func TestLimit2(t *testing.T) {
	src := []string{"a", "b", "c", "d", "e"}

	t.Run("limit 3", func(t *testing.T) {
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Limit2(3, slices.All(src))
			},
			Expected: []hiter.KeyValue[int, string]{{0, "a"}, {1, "b"}, {2, "c"}},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("limit 0", func(t *testing.T) {
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Limit2(0, slices.All(src))
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})
}

func TestMap(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}

	t.Run("square numbers", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Map(func(x int) int { return x * x }, slices.Values(src))
			},
			Expected: []int{1, 4, 9, 16, 25},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("int to string", func(t *testing.T) {
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.Map(func(x int) string { return fmt.Sprintf("num_%d", x) }, slices.Values(src))
			},
			Expected: []string{"num_1", "num_2", "num_3", "num_4", "num_5"},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("empty sequence", func(t *testing.T) {
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Map(func(x int) int { return x * 2 }, slices.Values([]int{}))
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})
}

func TestMap2(t *testing.T) {
	src := []string{"hello", "world", "test"}

	t.Run("key-value transformation", func(t *testing.T) {
		testcase.Two[string, int]{
			Seq: func() iter.Seq2[string, int] {
				return hiter.Map2(
					func(k int, v string) (string, int) { return v, k },
					slices.All(src),
				)
			},
			Expected: []hiter.KeyValue[string, int]{{"hello", 0}, {"world", 1}, {"test", 2}},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("string manipulation", func(t *testing.T) {
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Map2(
					func(k int, v string) (int, string) { return k * 10, strings.ToUpper(v) },
					slices.All(src),
				)
			},
			Expected: []hiter.KeyValue[int, string]{{0, "HELLO"}, {10, "WORLD"}, {20, "TEST"}},
			BreakAt:  2,
		}.Test(t)
	})
}

func TestReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		src := []int{1, 2, 3, 4, 5}
		result := hiter.Reduce(func(sum, v int) int { return sum + v }, 0, slices.Values(src))
		assert.Equal(t, 15, result)
	})

	t.Run("string concatenation", func(t *testing.T) {
		src := []string{"hello", " ", "world"}
		result := hiter.Reduce(func(acc, v string) string { return acc + v }, "", slices.Values(src))
		assert.Equal(t, "hello world", result)
	})

	t.Run("empty sequence", func(t *testing.T) {
		result := hiter.Reduce(func(sum, v int) int { return sum + v }, 42, slices.Values([]int{}))
		assert.Equal(t, 42, result)
	})

	t.Run("find maximum", func(t *testing.T) {
		src := []int{3, 1, 4, 1, 5, 9, 2, 6}
		result := hiter.Reduce(func(max, v int) int {
			if v > max {
				return v
			}
			return max
		}, src[0], slices.Values(src[1:]))
		assert.Equal(t, 9, result)
	})

	t.Run("build slice", func(t *testing.T) {
		src := []int{1, 2, 3}
		result := hiter.Reduce(func(acc []int, v int) []int {
			return append(acc, v*2)
		}, []int{}, slices.Values(src))
		assert.Assert(t, goCmp.DeepEqual([]int{2, 4, 6}, result))
	})
}

func TestReduce2(t *testing.T) {
	t.Run("sum with keys", func(t *testing.T) {
		src := []hiter.KeyValue[int, int]{{1, 10}, {2, 20}, {3, 30}}
		result := hiter.Reduce2(func(sum, k, v int) int { return sum + k + v }, 0, hiter.Values2(src))
		assert.Equal(t, 66, result) // 0 + (1+10) + (2+20) + (3+30) = 66
	})

	t.Run("build map", func(t *testing.T) {
		src := []hiter.KeyValue[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
		result := hiter.Reduce2(func(acc map[string]string, k string, v int) map[string]string {
			acc[k] = fmt.Sprintf("value_%d", v)
			return acc
		}, make(map[string]string), hiter.Values2(src))
		expected := map[string]string{"a": "value_1", "b": "value_2", "c": "value_3"}
		assert.Assert(t, goCmp.DeepEqual(expected, result))
	})

	t.Run("empty sequence", func(t *testing.T) {
		result := hiter.Reduce2(func(sum, k, v int) int { return sum + k + v }, 100, hiter.Values2([]hiter.KeyValue[int, int]{}))
		assert.Equal(t, 100, result)
	})

	t.Run("concatenate keys and values", func(t *testing.T) {
		src := []hiter.KeyValue[string, string]{{"hello", "world"}, {"foo", "bar"}}
		result := hiter.Reduce2(func(acc, k, v string) string {
			if acc == "" {
				return k + v
			}
			return acc + "_" + k + v
		}, "", hiter.Values2(src))
		assert.Equal(t, "helloworld_foobar", result)
	})
}

func TestZip(t *testing.T) {
	t.Run("equal length sequences", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 2, 3})
		seq2 := slices.Values([]string{"a", "b", "c"})

		var results []hiter.Zipped[int, string]
		for z := range hiter.Zip(seq1, seq2) {
			results = append(results, z)
		}

		expected := []hiter.Zipped[int, string]{
			{V1: 1, Ok1: true, V2: "a", Ok2: true},
			{V1: 2, Ok1: true, V2: "b", Ok2: true},
			{V1: 3, Ok1: true, V2: "c", Ok2: true},
		}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})

	t.Run("first sequence longer", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 2, 3, 4})
		seq2 := slices.Values([]string{"a", "b"})

		var results []hiter.Zipped[int, string]
		for z := range hiter.Zip(seq1, seq2) {
			results = append(results, z)
		}

		expected := []hiter.Zipped[int, string]{
			{V1: 1, Ok1: true, V2: "a", Ok2: true},
			{V1: 2, Ok1: true, V2: "b", Ok2: true},
			{V1: 3, Ok1: true, V2: "", Ok2: false},
			{V1: 4, Ok1: true, V2: "", Ok2: false},
		}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})

	t.Run("second sequence longer", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 2})
		seq2 := slices.Values([]string{"a", "b", "c", "d"})

		var results []hiter.Zipped[int, string]
		for z := range hiter.Zip(seq1, seq2) {
			results = append(results, z)
		}

		expected := []hiter.Zipped[int, string]{
			{V1: 1, Ok1: true, V2: "a", Ok2: true},
			{V1: 2, Ok1: true, V2: "b", Ok2: true},
			{V1: 0, Ok1: false, V2: "c", Ok2: true},
			{V1: 0, Ok1: false, V2: "d", Ok2: true},
		}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})

	t.Run("both empty sequences", func(t *testing.T) {
		seq1 := slices.Values([]int{})
		seq2 := slices.Values([]string{})

		var results []hiter.Zipped[int, string]
		for z := range hiter.Zip(seq1, seq2) {
			results = append(results, z)
		}

		assert.Assert(t, goCmp.Len(results, 0))
	})

	t.Run("early break", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 2, 3, 4, 5})
		seq2 := slices.Values([]string{"a", "b", "c", "d", "e"})

		var results []hiter.Zipped[int, string]
		for z := range hiter.Zip(seq1, seq2) {
			results = append(results, z)
			if len(results) == 2 {
				break
			}
		}

		expected := []hiter.Zipped[int, string]{
			{V1: 1, Ok1: true, V2: "a", Ok2: true},
			{V1: 2, Ok1: true, V2: "b", Ok2: true},
		}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})
}

func TestZip2(t *testing.T) {
	t.Run("equal length sequences", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{{0, "a"}, {1, "b"}})
		seq2 := hiter.Values2([]hiter.KeyValue[string, int]{{"x", 10}, {"y", 20}})

		var results []hiter.Zipped2[int, string, string, int]
		for z := range hiter.Zip2(seq1, seq2) {
			results = append(results, z)
		}

		expected := []hiter.Zipped2[int, string, string, int]{
			{K1: 0, V1: "a", Ok1: true, K2: "x", V2: 10, Ok2: true},
			{K1: 1, V1: "b", Ok1: true, K2: "y", V2: 20, Ok2: true},
		}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})

	t.Run("first sequence longer", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{{0, "a"}, {1, "b"}, {2, "c"}})
		seq2 := hiter.Values2([]hiter.KeyValue[string, int]{{"x", 10}})

		var results []hiter.Zipped2[int, string, string, int]
		for z := range hiter.Zip2(seq1, seq2) {
			results = append(results, z)
		}

		expected := []hiter.Zipped2[int, string, string, int]{
			{K1: 0, V1: "a", Ok1: true, K2: "x", V2: 10, Ok2: true},
			{K1: 1, V1: "b", Ok1: true, K2: "", V2: 0, Ok2: false},
			{K1: 2, V1: "c", Ok1: true, K2: "", V2: 0, Ok2: false},
		}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})

	t.Run("second sequence longer", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{{0, "a"}})
		seq2 := hiter.Values2([]hiter.KeyValue[string, int]{{"x", 10}, {"y", 20}, {"z", 30}})

		var results []hiter.Zipped2[int, string, string, int]
		for z := range hiter.Zip2(seq1, seq2) {
			results = append(results, z)
		}

		expected := []hiter.Zipped2[int, string, string, int]{
			{K1: 0, V1: "a", Ok1: true, K2: "x", V2: 10, Ok2: true},
			{K1: 0, V1: "", Ok1: false, K2: "y", V2: 20, Ok2: true},
			{K1: 0, V1: "", Ok1: false, K2: "z", V2: 30, Ok2: true},
		}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})

	t.Run("both empty sequences", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{})
		seq2 := hiter.Values2([]hiter.KeyValue[string, int]{})

		var results []hiter.Zipped2[int, string, string, int]
		for z := range hiter.Zip2(seq1, seq2) {
			results = append(results, z)
		}

		assert.Assert(t, goCmp.Len(results, 0))
	})

	t.Run("early break", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{{0, "a"}, {1, "b"}, {2, "c"}})
		seq2 := hiter.Values2([]hiter.KeyValue[string, int]{{"x", 10}, {"y", 20}, {"z", 30}})

		var results []hiter.Zipped2[int, string, string, int]
		for z := range hiter.Zip2(seq1, seq2) {
			results = append(results, z)
			if len(results) == 1 {
				break
			}
		}

		expected := []hiter.Zipped2[int, string, string, int]{{K1: 0, V1: "a", Ok1: true, K2: "x", V2: 10, Ok2: true}}
		assert.Assert(t, goCmp.DeepEqual(expected, results))
	})
}

func TestMerge(t *testing.T) {
	t.Run("merge sorted sequences", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 3, 5, 7})
		seq2 := slices.Values([]int{2, 4, 6, 8})
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Merge(seq1, seq2)
			},
			Expected: []int{1, 2, 3, 4, 5, 6, 7, 8},
			BreakAt:  4,
		}.Test(t)
	})

	t.Run("one empty sequence", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 2, 3})
		seq2 := slices.Values([]int{})
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Merge(seq1, seq2)
			},
			Expected: []int{1, 2, 3},
			BreakAt:  2,
		}.Test(t)
	})

	t.Run("both empty sequences", func(t *testing.T) {
		seq1 := slices.Values([]int{})
		seq2 := slices.Values([]int{})
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Merge(seq1, seq2)
			},
			Expected: nil,
			BreakAt:  0,
		}.Test(t)
	})

	t.Run("with duplicates", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 3, 5})
		seq2 := slices.Values([]int{1, 3, 4})
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.Merge(seq1, seq2)
			},
			Expected: []int{1, 1, 3, 3, 4, 5},
			BreakAt:  3,
		}.Test(t)
	})
}

func TestMergeFunc(t *testing.T) {
	t.Run("reverse order merge", func(t *testing.T) {
		seq1 := slices.Values([]int{7, 5, 3, 1})
		seq2 := slices.Values([]int{8, 6, 4, 2})
		testcase.One[int]{
			Seq: func() iter.Seq[int] {
				return hiter.MergeFunc(seq1, seq2, func(a, b int) int { return cmp.Compare(b, a) })
			},
			Expected: []int{8, 7, 6, 5, 4, 3, 2, 1},
			BreakAt:  4,
		}.Test(t)
	})

	t.Run("string length comparison", func(t *testing.T) {
		seq1 := slices.Values([]string{"a", "bbb", "eeeee"})
		seq2 := slices.Values([]string{"cc", "dddd"})
		testcase.One[string]{
			Seq: func() iter.Seq[string] {
				return hiter.MergeFunc(seq1, seq2, func(a, b string) int { return cmp.Compare(len(a), len(b)) })
			},
			Expected: []string{"a", "cc", "bbb", "dddd", "eeeee"},
			BreakAt:  3,
		}.Test(t)
	})
}

func TestMerge2(t *testing.T) {
	t.Run("merge sorted key-value sequences", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{{1, "a"}, {3, "c"}, {5, "e"}})
		seq2 := hiter.Values2([]hiter.KeyValue[int, string]{{2, "b"}, {4, "d"}, {6, "f"}})
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Merge2(seq1, seq2)
			},
			Expected: []hiter.KeyValue[int, string]{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}, {6, "f"}},
			BreakAt:  3,
		}.Test(t)
	})

	t.Run("one empty sequence", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{{1, "a"}, {2, "b"}})
		seq2 := hiter.Values2([]hiter.KeyValue[int, string]{})
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.Merge2(seq1, seq2)
			},
			Expected: []hiter.KeyValue[int, string]{{1, "a"}, {2, "b"}},
			BreakAt:  1,
		}.Test(t)
	})
}

func TestMergeFunc2(t *testing.T) {
	t.Run("merge by key comparison", func(t *testing.T) {
		seq1 := hiter.Values2([]hiter.KeyValue[int, string]{{10, "a"}, {30, "c"}})
		seq2 := hiter.Values2([]hiter.KeyValue[int, string]{{20, "b"}, {40, "d"}})
		testcase.Two[int, string]{
			Seq: func() iter.Seq2[int, string] {
				return hiter.MergeFunc2(seq1, seq2, func(k1, k2 int) int { return cmp.Compare(k1, k2) })
			},
			Expected: []hiter.KeyValue[int, string]{{10, "a"}, {20, "b"}, {30, "c"}, {40, "d"}},
			BreakAt:  2,
		}.Test(t)
	})
}
