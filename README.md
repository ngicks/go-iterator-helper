# go-iterator-helper

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/ngicks/go-iterator-helper)

Helpers / converters / sources for iterators.

## NOTE: things that are already iterator

Listed below are already iterators.
This module will not define iterator sources for these kind.
(in case older version of this modules has defined those, it should already be removed.)

```go
// https://pkg.go.dev/go/token#FileSet.Iterate
func (s *FileSet) Iterate(f func(*File) bool)
// https://pkg.go.dev/log/slog#Record.Attrs
func (r Record) Attrs(f func(Attr) bool)
// https://pkg.go.dev/sync#Map.Range
func (m *Map) Range(f func(key, value any) bool)
```

## Typical Usages

This section provides examples of common iterator operations. Each example has a dedicated example test.

### 1. Basic Operations

#### Creating Iterators

See [`Example_creatingIterators`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-CreatingIterators) for creating iterators from slices, ranges, repeated values, and single values.

```go
// From slices, Range, Repeat, Once
// From slice (std): 1 - 3
slice_iter := slices.Values([]int{1, 2, 3})
// Range: 10 - 14 (inclusive-exclusive range)
// To control inclusiveness, use RangeInclusive
range_iter := hiter.Range(10, 15)
// Repeat: [hi hi hi]
repeat_iter := hiter.Repeat("hi", 3)
// Once: [42]
once_iter := hiter.Once(42)
```

#### Combining Iterators

See [`Example_concat`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-Concat) for combining multiple iterators.

```go
first := slices.Values([]int{1, 2, 3})
second := slices.Values([]int{10, 11})
third := slices.Values([]int{20, 21, 22})

combined := hiter.Concat(first, second, third)
// Output: [1 2 3 10 11 20 21 22]
```

#### Transforming Data

See [`Example_mapAndFilter`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-MapAndFilter) for transforming iterators with Map and Filter.

```go
numbers := slices.Values([]int{1, 2, 3, 4, 5, 6})

doubled := hiter.Map(func(n int) int { return n * 2 }, numbers)
// 2, 4, 6, 8, 10, 12
evens := hiter.Filter(func(n int) bool { return n%2 == 0 }, numbers)
// 2, 4, 6
```

#### Aggregation

See [`Example_reduce`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-Reduce) for basic aggregation operations.

```go
numbers := slices.Values([]int{1, 2, 3, 4, 5})

sum := hiter.Reduce(func(acc, val int) int { return acc + val }, 0, numbers)
// Sum: 15
// For this simple case use instead
sum := hiter.Sum(numbers)
// Sum: 15

product := hiter.Reduce(func(acc, val int) int { return acc * val }, 1, numbers)
// Product: 120
```

#### Group and Reduce

See [`Example_reduceGroup`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-ReduceGroup) for grouping by key and aggregating values.

```go
grouped := hiter.ReduceGroup(
    func(acc, val int) int { return acc + val },
    0,
    pairs,
)
// Result: map[fruit:55 vegetable:45]
```

#### String Operations

See [`Example_stringJoin`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-StringJoin) for joining iterator values into strings.

```go
words := slices.Values([]string{"go", "iterator", "helper"})
collected := stringsiter.Collect(words)
// Result: "goiteratorhelper"
joined := stringsiter.Join("-", words)
// Result: "go-iterator-helper"
```

### 2. Error Handling

#### Error-Aware Iteration

See [`Example_tryForEach`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-TryForEach) for processing until errors occur.

```go
err := hiter.TryForEach(func(item string) { /* process */ }, data)
```

#### JSON Stream Error Handling

See [`Example_errboxJSON`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-ErrboxJSON) for handling errors in JSON streams.

```go
decoder := json.NewDecoder(reader)
jsonBox := errbox.New(encodingiter.Decode[Person](decoder))
for person := range jsonBox.IntoIter() {
    // Process valid records
}
if err := jsonBox.Err(); err != nil {
    // Handle stream errors
}
```

### 3. Advanced Patterns

#### Moving Averages with Windows

See [`Example_windowMovingAverage`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-WindowMovingAverage) for sliding window calculations.

```go
windows := hiter.Window(data, 3)
averages := hiter.Map(func(w []int) float64 {
    return float64(hiter.Sum(slices.Values(w))) / float64(len(w))
}, windows)
```

#### Flattening Nested Structures

See [`Example_flatten`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-Flatten) for flattening nested arrays.

```go
flattened := hiter.Flatten(slices.Values(nested))
// [[a b] [c] [d e f]] → [a b c d e f]
```

#### Resumable Iteration

See [`Example_resumable`](https://pkg.go.dev/github.com/ngicks/go-iterator-helper#example-package-Resumable) for pausable and resumable iteration.

```go
resumable := iterable.NewResumable(source)
batch1 := hiter.Limit(3, resumable.IntoIter())  // [1 2 3]
batch2 := hiter.Limit(3, resumable.IntoIter())  // [4 5 6]
remaining := resumable.IntoIter()               // [7 8 9]
```

## hiter

Helpers for iterator.

This package avoids re-implementing those which defined in standard or quasi-standard libraries.
Namely `slices`, `maps`, ~~`x/exp/xiter`~~ (the proposal was withdrawn, at this time I'm not seeing what will be implemented in the furure hypothesis `x/exp/xiter`).

~~For example, `Zip`, `Reduce` are not defined since they will be implemented in `xiter` when~~
~~[#61898](https://github.com/golang/go/issues/61898) accepted and merged.~~

Some ideas are stolen from https://jsr.io/@std/collections/doc, like Permutation and SumOf.

Each package named `*iter` corresponds to same `*` name of std library (e.g. `reflectiter` defines iterator souces/adapters for std package `refect`).
Packages nested under other package are flattened, e.g. `encodingiter` defines helpers for `encoding/json`, `encoding/xml`, `encoding/csv` and so on.

## Deprecated: x/exp/xiter

Deprecated: you should no longer use this package since the proposal is withdrawn.
`hiter` re-defines equivalents so you can use these in there.
The proposal was wound down because the author saw Go iterator was too young.
Once Go iterator gets matured in the community, proposal might again be proposed.
At that time signatures of functions would be changed if the community finds better conventions.

Those listed in [#61898](https://github.com/golang/go/issues/61898).

This package is vendored so that you can use it anywhere without copy-and-pasting everywhere.
It is already frozen; no change will be made even when xiter proposal got some modification.

## Future deprecations

All functions will be noted as deprecated when std, `golang.org/x/exp/xiter` or similar quasi-std packages define equivalents.

It is just simply noted as deprecated: functions will remain same regardless of deprecation.
You definitely should use std where possible, but you can keep using `hiter`.

## Deprecated functions

### After Go 1.24

Nothing.

https://tip.golang.org/doc/go1.24

`Go 1.24` adds

- `Line`, `SplitSeq`, `SplitAfterSeq`, `FieldSeq` and `FieldFuncSeq` to `(strings|bytes)`.
  - `stringsiter` package defines similar functions but not exactly same. Those remain valid and maintained.
- iterator sources to `go/types`
  - `hiter.Atter` remains valid since it can be used for any implementors, e.g. [github.com/gammazero/deque](https://github.com/gammazero/deque)

### After Go 1.25

Seems nothing

https://tip.golang.org/doc/go1.25
