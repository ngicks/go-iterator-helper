# go-iterator-helper

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/ngicks/go-iterator-helper)

Helpers / converters / sources for iterators.

## hiter

Helpers for iterator.

This package avoids re-implementing those which defined in standard or quasi-standard libraries.
Namely `slices`, `maps`, `x/exp/xiter`.

For example, `Zip`, `Reduce` are not defined since they will be implemented in `xiter` when
[#61898](https://github.com/golang/go/issues/61898) accepted and merged.

Some ideas are stolen from https://jsr.io/@std/collections/doc, like Permutation and SumOf.

Iterator sources: functions that compose up iterators from data sources:

```go
func Chan[V any](ctx context.Context, ch <-chan V) iter.Seq[V]
func Heap[V any](h heap.Interface) iter.Seq[V]
func IndexAccessible[A Atter[T], T any](a A, indices iter.Seq[int]) iter.Seq2[int, T]
func JsonDecoder(dec *json.Decoder) iter.Seq2[json.Token, error]
func ListAll[V any](l *list.List) iter.Seq[V]
func ListBackward[V any](l *list.List) iter.Seq[V]
func ListElementAll[V any](ele *list.Element) iter.Seq[V]
func ListElementBackward[V any](ele *list.Element) iter.Seq[V]
func MergeSort[S ~[]T, T cmp.Ordered](m S) iter.Seq[T]
func MergeSortFunc[S ~[]T, T any](m S, cmp func(l, r T) int) iter.Seq[T]
func MergeSortSliceLike[S SliceLike[T], T cmp.Ordered](s S) iter.Seq[T]
func MergeSortSliceLikeFunc[S SliceLike[T], T any](s S, cmp func(l, r T) int) iter.Seq[T]
func Permutations[S ~[]E, E any](in S) iter.Seq[S]
func Range[T Numeric](start, end T) iter.Seq[T]
func Repeat[V any](v V, n int) iter.Seq[V]
func Repeat2[K, V any](k K, v V, n int) iter.Seq2[K, V]
func RepeatFunc[V any](fnV func() V, n int) iter.Seq[V]
func RepeatFunc2[K, V any](fnK func() K, fnV func() V, n int) iter.Seq2[K, V]
func RingAll[V any](r *ring.Ring) iter.Seq[V]
func RingBackward[V any](r *ring.Ring) iter.Seq[V]
func RunningReduce[V, Sum any](reducer func(accumulator Sum, current V, i int) Sum, initial Sum, ...) iter.Seq[Sum]
func Scan(scanner *bufio.Scanner) iter.Seq[string]
func SqlRows[T any](r *sql.Rows, scanner func(*sql.Rows) (T, error)) iter.Seq2[T, error]
func StringsChunk(s string, n int) iter.Seq[string]
func StringsRuneChunk(s string, n int) iter.Seq[string]
func StringsSplitFunc(s string, n int, splitFn StringsCutterFunc) iter.Seq[string]
func SyncMap[K, V any](m *sync.Map) iter.Seq2[K, V]
func Window[S ~[]E, E any](s S, n int) iter.Seq[S]
func XmlDecoder(dec *xml.Decoder) iter.Seq2[xml.Token, error]
```

Iterator adapters: iterator that processes / modifies values from other iterators.

```go
func Alternate[V any](seqs ...iter.Seq[V]) iter.Seq[V]
func Alternate2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V]
func Compact[V comparable](seq iter.Seq[V]) iter.Seq[V]
func Compact2[K, V comparable](seq iter.Seq2[K, V]) iter.Seq2[K, V]
func CompactFunc[V any](eq func(i, j V) bool, seq iter.Seq[V]) iter.Seq[V]
func CompactFunc2[K, V any](eq func(k1 K, v1 V, k2 K, v2 V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V]
func Decorate[V any](prepend, append Iterable[V], seq iter.Seq[V]) iter.Seq[V]
func Decorate2[K, V any](prepend, append Iterable2[K, V], seq iter.Seq2[K, V]) iter.Seq2[K, V]
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T]
func Flatten[S ~[]E, E any](seq iter.Seq[S]) iter.Seq[E]
func FlattenF[S1 ~[]E1, E1 any, E2 any](seq iter.Seq2[S1, E2]) iter.Seq2[E1, E2]
func FlattenL[S2 ~[]E2, E1 any, E2 any](seq iter.Seq2[E1, S2]) iter.Seq2[E1, E2]
func LimitUntil[V any](f func(V) bool, seq iter.Seq[V]) iter.Seq[V]
func LimitUntil2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V]
func Omit[K any](seq iter.Seq[K]) func(yield func() bool)
func Omit2[K, V any](seq iter.Seq2[K, V]) func(yield func() bool)
func OmitF[K, V any](seq iter.Seq2[K, V]) iter.Seq[V]
func OmitL[K, V any](seq iter.Seq2[K, V]) iter.Seq[K]
func Pairs[K, V any](seq1 iter.Seq[K], seq2 iter.Seq[V]) iter.Seq2[K, V]
func Skip[V any](n int, seq iter.Seq[V]) iter.Seq[V]
func Skip2[K, V any](n int, seq iter.Seq2[K, V]) iter.Seq2[K, V]
func SkipLast[V any](n int, seq iter.Seq[V]) iter.Seq[V]
func SkipLast2[K, V any](n int, seq iter.Seq2[K, V]) iter.Seq2[K, V]
func SkipWhile[V any](f func(V) bool, seq iter.Seq[V]) iter.Seq[V]
func SkipWhile2[K, V any](f func(K, V) bool, seq iter.Seq2[K, V]) iter.Seq2[K, V]
func Tap[V any](tap func(V), seq iter.Seq[V]) iter.Seq[V]
func Tap2[K, V any](tap func(K, V), seq iter.Seq2[K, V]) iter.Seq2[K, V]
func Transpose[K, V any](seq iter.Seq2[K, V]) iter.Seq2[V, K]
func WindowSeq[V any](n int, seq iter.Seq[V]) iter.Seq[iter.Seq[V]]
```

Collectors: functions that collect data from iterators and convert to other data.

```go
func AppendSeq2[S ~[]KeyValue[K, V], K, V any](s S, seq iter.Seq2[K, V]) S
func ChanSend[V any](ctx context.Context, c chan<- V, seq iter.Seq[V]) (v V, sentAll bool)
func Find[V comparable](v V, seq iter.Seq[V]) (V, int)
func Find2[K, V comparable](k K, v V, seq iter.Seq2[K, V]) (K, V, int)
func FindFunc[V any](f func(V) bool, seq iter.Seq[V]) (V, int)
func FindFunc2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) (K, V, int)
func FindLast[V comparable](v V, seq iter.Seq[V]) (found V, idx int)
func FindLast2[K, V comparable](k K, v V, seq iter.Seq2[K, V]) (foundK K, foundV V, idx int)
func FindLastFunc[V any](fn func(V) bool, seq iter.Seq[V]) (found V, idx int)
func FindLastFunc2[K, V any](fn func(K, V) bool, seq iter.Seq2[K, V]) (foundK K, foundV V, idx int)
func ReduceGroup[K comparable, V, Sum any](reducer func(accumulator Sum, current V) Sum, initial Sum, seq iter.Seq2[K, V]) map[K]Sum
func StringsCollect(sizeHint int, seq iter.Seq[string]) string
func Sum[S Summable](seq iter.Seq[S]) S
func SumOf[V any, S Summable](selector func(ele V) S, seq iter.Seq[V]) S
func Collect2[K, V any](seq iter.Seq2[K, V]) []KeyValue[K, V]
```

## hiter/iterable

Wrapper for iterable objects; heap, list, ring, slice, map, channel, etc.

All of them implement 1 or 2 of `Iter() iter.Seq[V]`, `Iter2() iter.Seq[K, V]`, `IntoIter() iter.Seq[V]` or `IntoIter2() iter.Seq2[K, V]`

```go
package iterable // import "github.com/ngicks/go-iterator-helper/hiter/iterable"

type Chan[V any] struct{ ... }
type Heap[T any] struct{ ... }
type IndexAccessible[A hiter.Atter[T], T any] struct{ ... }
type JsonDecoder struct{ ... }
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
type RingAll[T any] struct{ ... }
type RingBackward[T any] struct{ ... }
type Scanner struct{ ... }
type SliceAll[E any] []E
type SliceBackward[E any] []E
type SqlRows[T any] struct{ ... }
type SyncMap[K comparable, V any] struct{ ... }
type XmlDecoder struct{ ... }

```

## hiter/errbox

`hiter/errbox` defines an utility that wraps `iter.Seq2[V, error]` to `iter.Seq[V]` by remembering first error encountered.

```go
package errbox // import "github.com/ngicks/go-iterator-helper/hiter/errbox"

type Box[V any] struct{ ... }
    func New[V any](seq iter.Seq2[V, error]) *Box[V]
type JsonDecoder struct{ ... }
    func NewJsonDecoder(dec *json.Decoder) *JsonDecoder
type SqlRows[V any] struct{ ... }
    func NewSqlRows[V any](rows *sql.Rows, scanner func(*sql.Rows) (V, error)) *SqlRows[V]
type XmlDecoder struct{ ... }
    func NewXmlDecoder(dec *xml.Decoder) *XmlDecoder

```

## x/exp/xiter

Those listed in [#61898](https://github.com/golang/go/issues/61898).

This package is vendored so that you can use it anywhere without copy-and-pasting everywhere.
It is already frozen; no change will be made even when xiter proposal got some modification.
