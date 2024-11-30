package hiter_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/internal/testcase"
)

func TestReplace(t *testing.T) {
	seq := slices.Values([]int{0, 5, 1, 2, 5, 6, 5})
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return hiter.Replace(5, -1, -1, seq)
		},
		Expected: []int{0, -1, 1, 2, -1, 6, -1},
		BreakAt:  2,
	}.Test(t)
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return hiter.Replace(5, -1, 2, seq)
		},
		Expected: []int{0, -1, 1, 2, -1, 6, 5},
		BreakAt:  2,
	}.Test(t)
}

func TestReplaceFunc(t *testing.T) {
	seq := slices.Values([]int{0, 5, 1, 2, 5, 6, 5})
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return hiter.ReplaceFunc(
				func(i int) bool { return i >= 5 },
				-1,
				-1,
				seq,
			)
		},
		Expected: []int{0, -1, 1, 2, -1, -1, -1},
		BreakAt:  2,
	}.Test(t)
	testcase.One[int]{
		Seq: func() iter.Seq[int] {
			return hiter.ReplaceFunc(
				func(i int) bool { return i >= 5 },
				-1,
				3,
				seq,
			)
		},
		Expected: []int{0, -1, 1, 2, -1, -1, 5},
		BreakAt:  2,
	}.Test(t)
}

func TestReplace2(t *testing.T) {
	seq := hiter.Values2([]hiter.KeyValue[int, int]{
		{K: 1, V: 0},
		{K: 5, V: 2},
		{K: 1, V: 1},
		{K: 2, V: 4},
		{K: 5, V: 2},
		{K: 6, V: 2},
		{K: 5, V: 2},
	})
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Replace2(5, 2, -1, -2, -1, seq)
		},
		Expected: []hiter.KeyValue[int, int]{
			{K: 1, V: 0},
			{K: -1, V: -2},
			{K: 1, V: 1},
			{K: 2, V: 4},
			{K: -1, V: -2},
			{K: 6, V: 2},
			{K: -1, V: -2},
		},
		BreakAt: 2,
	}.Test(t)
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.Replace2(5, 2, -1, -2, 2, seq)
		},
		Expected: []hiter.KeyValue[int, int]{
			{K: 1, V: 0},
			{K: -1, V: -2},
			{K: 1, V: 1},
			{K: 2, V: 4},
			{K: -1, V: -2},
			{K: 6, V: 2},
			{K: 5, V: 2},
		},
		BreakAt: 2,
	}.Test(t)
}

func TestReplaceFunc2(t *testing.T) {
	seq := hiter.Values2([]hiter.KeyValue[int, int]{
		{K: 1, V: 0},
		{K: 5, V: 2},
		{K: 1, V: 1},
		{K: 2, V: 4},
		{K: 5, V: 2},
		{K: 6, V: 2},
		{K: 5, V: 2},
	})
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.ReplaceFunc2(
				func(k, v int) bool {
					return k == 2 || v == 2
				},
				-1, -2,
				-1,
				seq,
			)
		},
		Expected: []hiter.KeyValue[int, int]{
			{K: 1, V: 0},
			{K: -1, V: -2},
			{K: 1, V: 1},
			{K: -1, V: -2},
			{K: -1, V: -2},
			{K: -1, V: -2},
			{K: -1, V: -2},
		},
		BreakAt: 2,
	}.Test(t)
	testcase.Two[int, int]{
		Seq: func() iter.Seq2[int, int] {
			return hiter.ReplaceFunc2(
				func(k, v int) bool {
					return k == 2 || v == 2
				},
				-1, -2,
				2,
				seq,
			)
		},
		Expected: []hiter.KeyValue[int, int]{
			{K: 1, V: 0},
			{K: -1, V: -2},
			{K: 1, V: 1},
			{K: -1, V: -2},
			{K: 5, V: 2},
			{K: 6, V: 2},
			{K: 5, V: 2},
		},
		BreakAt: 2,
	}.Test(t)
}
