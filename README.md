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

## hiter

Helpers for iterator.

This package avoids re-implementing those which defined in standard or quasi-standard libraries.
Namely `slices`, `maps`, ~~`x/exp/xiter`~~ (the proposal was withdrawn, at this time I'm not seeing what will be implemented in `x/exp/xiter`).

For example, `Zip`, `Reduce` are not defined since they will be implemented in `xiter` when
[#61898](https://github.com/golang/go/issues/61898) accepted and merged.

Some ideas are stolen from https://jsr.io/@std/collections/doc, like Permutation and SumOf.

Each package named `*iter` corresponds to same `*` name of std library (e.g. `reflectiter` defines iterator souces/adapters for std package `refect`).
Packages nested under other package are flattened, e.g. `encodingiter` defines helpers for `encoding/json`, `encoding/xml`, `encoding/csv` and so on.

## x/exp/xiter

Those listed in [#61898](https://github.com/golang/go/issues/61898).

Deprecated: you should no longer use this package since the proposal is withdrawn.
`hiter` re-defines equivalents so you can use these in there.
The proposal was wound down because the author saw Go iterator was too young.
Once Go iterator gets matured in the community, proposal might again be proposed.
At that time signatures of functions would be changed if the community finds better conventions.

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
