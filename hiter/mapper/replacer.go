package mapper

import (
	"iter"
	"strings"

	"github.com/ngicks/go-iterator-helper/hiter"
)

type StringsReplacer struct {
	r *strings.Replacer
}

func NewStringsReplacer(oldnew ...string) *StringsReplacer {
	return &StringsReplacer{strings.NewReplacer(oldnew...)}
}

func (r *StringsReplacer) Map(seq iter.Seq[string]) iter.Seq[string] {
	return hiter.Map(r.r.Replace, seq)
}

type Replacer[K comparable, V any] map[K]V

func (r Replacer[K, V]) Map(seq iter.Seq[K]) iter.Seq[V] {
	return hiter.Map(func(k K) V { return r[k] }, seq)
}
