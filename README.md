# go-iterator-helper

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/ngicks/go-iterator-helper)

Helpers / converters / sources for iterators.

## NOTE: things that are already iterator

Listed below are already iterator.
This module will not defines iterator sources for these kind.
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
Namely `slices`, `maps`, `x/exp/xiter`.

For example, `Zip`, `Reduce` are not defined since they will be implemented in `xiter` when
[#61898](https://github.com/golang/go/issues/61898) accepted and merged.

Some ideas are stolen from https://jsr.io/@std/collections/doc, like Permutation and SumOf.

## x/exp/xiter

Those listed in [#61898](https://github.com/golang/go/issues/61898).

This package is vendored so that you can use it anywhere without copy-and-pasting everywhere.
It is already frozen; no change will be made even when xiter proposal got some modification.
