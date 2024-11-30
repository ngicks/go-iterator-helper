/*
Wrapper for iterable objects; heap, list, ring, slice, map, channel, etc.
All of them implement 1 or 2 of `Iter() iter.Seq[V]`, `Iter2() iter.Seq[K, V]`, `IntoIter() iter.Seq[V]` or `IntoIter2() iter.Seq2[K, V]`
*/
package iterable
