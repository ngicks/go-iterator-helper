# go-iterator-helper

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/ngicks/go-iterator-helper)

Helpers / converters / sources for iterators.

## versioning policy

Before 1.0.0, sometimes packages might break anything in the module.
The breakage should be kept as little as possible between tags.

Every untagged commits break added/modified things.

## hiter

Helpers for iterator.

This package avoids re-implementing those which defined in standard or quasi-standard libraries.
Namely `slices`, `maps`, `x/exp/xiter`.

For example, `Zip`, `Reduce` are not defined since they will be implemented in `xiter` when
[#61898](https://github.com/golang/go/issues/61898) accepted and merged.

Some ideas are stolen from https://jsr.io/@std/collections/doc, like Permutation and SumOf.

Iterator sources: functions that compose up iterators from data sources:

```go
{{join .GoDoc.Hiter.Source "\n"}}
```

Iterator adapters: iterator that processes / modifies values from other iterators.

```go
{{join .GoDoc.Hiter.Adapter "\n"}}
```

Collectors: functions that collect data from iterators and convert to other data.

```go
{{join .GoDoc.Hiter.Collector "\n"}}
```

## hiter/iterable

Wrapper for iterable objects; heap, list, ring, slice, map, channel, etc.

All of them implement 1 or 2 of `Iter() iter.Seq[V]`, `Iter2() iter.Seq[K, V]`, `IntoIter() iter.Seq[V]` or `IntoIter2() iter.Seq2[K, V]`

```go
{{.GoDoc.HiterIterable}}
```

## hiter/errbox

`hiter/errbox` defines an utility that wraps `iter.Seq2[V, error]` to `iter.Seq[V]` by remembering the first error encountered.

```go
{{.GoDoc.HiterErrbox}}
```

## hiter/async

`hiter/async` defines asynchronous adapters

```go
{{.GoDoc.HiterAsync}}
```

## hiter/sh

Some short hands for adapters.
These are implemented only combining other components defined in this module (including `x/exp/xiter`).

```go
{{.GoDoc.HiterSh}}
```

## x/exp/xiter

Those listed in [#61898](https://github.com/golang/go/issues/61898).

This package is vendored so that you can use it anywhere without copy-and-pasting everywhere.
It is already frozen; no change will be made even when xiter proposal got some modification.
