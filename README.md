# go-iterator-helper

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/ngicks/go-iterator-helper)

Helpers / converters / sources for iterators.

## collection

Some useful function consuming iter.Seq / iter.Seq2.

The idea is stolen from https://jsr.io/@std/collections/doc.

```go
func Permutations[S ~[]E, E any](in S) iter.Seq[S]
func ReduceGroup[K comparable, V, Sum any](seq iter.Seq2[K, V], reducer func(accumulator Sum, current V) Sum, initial Sum) map[K]Sum
func RunningReduce[V, Sum any](seq iter.Seq[V], reducer func(accumulator Sum, current V, i int) Sum, ...) iter.Seq[Sum]
func SumOf[T any, E hiter.Numeric](seq iter.Seq[T], selector func(ele T) E) E
```

## hiter

Helpers for iterator.

This package avoids re-implementing those who defined in standard or quasi-standard libraries.
Namely `slices`, `maps`, `x/exp/xiter`.

For example, `Zip`, `Reduce` are not defined since it will be implemented in `xiter` when
[#61898](https://github.com/golang/go/issues/61898) accepted and merged.

```go
func Alternate[V any](seqs ...iter.Seq[V]) iter.Seq[V]
func Alternate2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V]
func AppendSeq2[S ~[]KeyValue[K, V], K, V any](s S, seq iter.Seq2[K, V]) S
func Chan[V any](ch <-chan V, f func()) iter.Seq[V]
func Decorate[V any](seq iter.Seq[V], prepend, append Iterable[V]) iter.Seq[V]
func Decorate2[K, V any](seq iter.Seq2[K, V], prepend, append Iterable2[K, V]) iter.Seq2[K, V]
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T]
func Flatten[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[E]
func FlattenF[S1 ~[]E1, E1 any, E2 any](i iter.Seq2[S1, E2]) iter.Seq2[E1, E2]
func FlattenL[S2 ~[]E2, E1 any, E2 any](i iter.Seq2[E1, S2]) iter.Seq2[E1, E2]
func Heap[T any](h heap.Interface) iter.Seq[T]
func LimitUntil[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V]
func LimitUntil2[K, V any](seq iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V]
func ListAll[T any](l *list.List) iter.Seq[T]
func ListBackward[T any](l *list.List) iter.Seq[T]
func ListElementAll[T any](ele *list.Element) iter.Seq[T]
func ListElementBackward[T any](ele *list.Element) iter.Seq[T]
func Omit[K any](seq iter.Seq[K]) func(yield func() bool)
func Omit2[K, V any](seq iter.Seq2[K, V]) func(yield func() bool)
func OmitF[T, U any](i iter.Seq2[T, U]) iter.Seq[U]
func OmitL[T, U any](i iter.Seq2[T, U]) iter.Seq[T]
func Pairs[K, V any](seq1 iter.Seq[K], seq2 iter.Seq[V]) iter.Seq2[K, V]
func Range[T Numeric](start, end T) iter.Seq[T]
func Repeat[V any](v V, n int) iter.Seq[V]
func Repeat2[K, V any](k K, v V, n int) iter.Seq2[K, V]
func RepeatFunc[V any](fnV func() V, n int) iter.Seq[V]
func RepeatFunc2[K, V any](fnK func() K, fnV func() V, n int) iter.Seq2[K, V]
func RingAll[T any](r *ring.Ring) iter.Seq[T]
func RingBackward[T any](r *ring.Ring) iter.Seq[T]
func Scan(scanner *bufio.Scanner) iter.Seq2[string, error]
func Skip[V any](seq iter.Seq[V], n int) iter.Seq[V]
func Skip2[K, V any](seq iter.Seq2[K, V], n int) iter.Seq2[K, V]
func SkipLast[V any](seq iter.Seq[V], n int) iter.Seq[V]
func SkipLast2[K, V any](seq iter.Seq2[K, V], n int) iter.Seq2[K, V]
func SkipWhile[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V]
func SkipWhile2[K, V any](seq iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V]
func StringsChunk(s string, n int) iter.Seq[string]
func StringsCollect(seq iter.Seq[string], sizeHint int) string
func StringsCutNewLine(s string) (int, int)
func StringsCutUpperCase(s string) (tokUntil int, skipUntil int)
func StringsCutWord(s string) (tokUntil int, skipUntil int)
func StringsRuneChunk(s string, n int) iter.Seq[string]
func StringsSplitFunc(s string, n int, splitFn StringsCutterFunc) iter.Seq[string]
func SyncMap[K, V any](m *sync.Map) iter.Seq2[K, V]
func Transpose[K, V any](seq iter.Seq2[K, V]) iter.Seq2[V, K]
func Window[S ~[]E, E any](s S, n int) iter.Seq[S]
type FuncIterable[V any] func() iter.Seq[V]
type FuncIterable2[K, V any] func() iter.Seq2[K, V]
type IntoIterable[V any] interface{ ... }
type IntoIterable2[K, V any] interface{ ... }
type Iterable[V any] interface{ ... }
type Iterable2[K, V any] interface{ ... }
type KeyValue[K, V any] struct{ ... }
    func Collect2[K, V any](seq iter.Seq2[K, V]) []KeyValue[K, V]
type KeyValues[K, V any] []KeyValue[K, V]
type Numeric interface{ ... }
type StringsCutterFunc func(s string) (tokUntil, skipUntil int)
```

## hiter/iterable

Wrapper for iterable objects; heap, list, ring, slice, map, channel, etc.

All of them implement 1 or 2 of `Iter() iter.Seq[V]`, `Iter2() iter.Seq[K, V]`, `IntoIter() iter.Seq[V]` or `IntoIter2() iter.Seq2[K, V]`

```go
type Chan[V any] <-chan V
type Heap[T any] struct{ ... }
type ListAll[T any] struct{ ... }
type ListBackward[T any] struct{ ... }
type ListElementAll[T any] struct{ ... }
type ListElementBackward[T any] struct{ ... }
type MapAll[K comparable, V any] map[K]V
type MapSorted[K cmp.Ordered, V any] map[K]V
type MapSortedFunc[M ~map[K]V, K comparable, V any] struct{ ... }
type Range[T hiter.Numeric] struct{ ... }
type Repeatable[V any] struct{ ... }
type Repeatable2[K, V any] struct{ ... }
type RepeatableFunc[V any] struct{ ... }
type RepeatableFunc2[K, V any] struct{ ... }
type Ring[T any] struct{ ... }
type RingBackward[T any] struct{ ... }
type Scanner struct{ ... }
type SliceAll[E any] []E
type SliceBackward[E any] []E
type SyncMap[K comparable, V any] struct{ ... }
```

## x/exp/xiter

Those listed in [#61898](https://github.com/golang/go/issues/61898).

This package is vendored so that you can use it anywhere without copy-and-pasting everywhere.
It is already frozen; no change will be made even when xiter proposal got some modification.
