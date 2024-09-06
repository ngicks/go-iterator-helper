# go-iterator-helper

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/ngicks/go-iterator-helper)

Helpers / converters / sources for iterators.

## collection

Some useful function consuming iter.Seq / iter.Seq2.

The idea is stolen from https://jsr.io/@std/collections/doc.

```go
{{.GoDoc.Collection}}
```

## hiter

Helpers for iterator.

This package avoids re-implementing those who defined in standard or quasi-standard libraries.
Namely `slices`, `maps`, `x/exp/xiter`.

For example, `Zip`, `Reduce` are not defined since it will be implemented in `xiter` when
[#61898](https://github.com/golang/go/issues/61898) accepted and merged.

```go
{{.GoDoc.Hiter}}
```

## hiter/iterable

Wrapper for iterable objects; heap, list, ring, slice, map, channel, etc.

All of them implement 1 or 2 of `Iter() iter.Seq[V]`, `Iter2() iter.Seq[K, V]`, `IntoIter() iter.Seq[V]` or `IntoIter2() iter.Seq2[K, V]`

```go
{{.GoDoc.HiterIterable}}
```

## x/exp/xiter

Those listed in [#61898](https://github.com/golang/go/issues/61898).

This package is vendored so that you can use it anywhere without copy-and-pasting everywhere.
It is already frozen; no change will be made even when xiter proposal got some modification.
