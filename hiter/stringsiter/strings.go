package stringsiter

import (
	"iter"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ngicks/go-iterator-helper/hiter"
)

// Collect concatenates all values form seq to a single string.
// sizeHint hints size of all values in bytes,
// which will be used to pre-allocate buffer.
func Collect(sizeHint int, seq iter.Seq[string]) string {
	var buf strings.Builder
	buf.Grow(sizeHint)
	for s := range seq {
		buf.WriteString(s)
	}
	return buf.String()
}

// Join is like [Collect] but inserts sep between every 2 values from seq, corresponds to [strings.Join].
func Join(sizeHint int, sep string, seq iter.Seq[string]) string {
	return Collect(
		sizeHint,
		hiter.Skip(
			1,
			hiter.Decorate(
				hiter.WrapSeqIterable(hiter.Once(sep)),
				nil,
				seq,
			),
		),
	)
}

// Chunk returns an iterator over non overlapping sub strings of n bytes.
// Sub slicing may cut in mid of utf8 sequences.
func Chunk(s string, n int) iter.Seq[string] {
	return func(yield func(string) bool) {
		s := s // no state in the seq.
		if n <= 0 {
			return
		}
		var cut string
		for {
			if len(s) >= n {
				cut, s = s[:n], s[n:]
			} else {
				cut, s = s, ""
			}
			if cut == "" {
				return
			}
			if !yield(cut) {
				return
			}
		}
	}
}

// RuneChunk returns an iterator over non overlapping sub strings of n utf8 characters.
func RuneChunk(s string, n int) iter.Seq[string] {
	return func(yield func(string) bool) {
		s := s // no state in the seq.
		for len(s) > 0 {
			var i int
			for range n {
				_, j := utf8.DecodeRuneInString(s[i:])
				if j == 0 {
					break
				}
				i += j
			}
			if i == 0 {
				return
			}
			if !yield(s[:i]) {
				return
			}
			s = s[i:]
		}
	}
}

// CutterFunc is used with [SplitFunc] to cut string from head.
// s[:tokUntil] is yielded through StringsSplitFunc.
// s[tokUntil:skipUntil] will be ignored.
type CutterFunc func(s string) (tokUntil, skipUntil int)

// CutNewLine is used with [SplitFunc].
// The input strings will be splitted at "\n".
// It also skips "\r" preceding "\n".
func CutNewLine(s string) (tokUntil int, skipUntil int) {
	i := strings.Index(s, "\n")
	j := i + 1
	if i >= 1 && strings.HasPrefix(s[i-1:], "\r\n") {
		i--
	}
	return i, j
}

// CutWord is a split function for a [SplitFunc] that returns each space-separated word of text,
// with surrounding spaces deleted. It will never return an empty string.
// The definition of space is set by unicode.IsSpace.
func CutWord(s string) (tokUntil int, skipUntil int) {
	if len(s) < 1 {
		return len(s), len(s)
	}
	var i int
	for len(s) > 0 {
		r, k := utf8.DecodeRuneInString(s)
		if unicode.IsSpace(r) {
			break
		}
		s = s[k:]
		i += k
	}
	j := i
	for len(s) > 0 {
		r, k := utf8.DecodeRuneInString(s)
		if !unicode.IsSpace(r) {
			break
		}
		s = s[k:]
		j += k
	}
	return i, j
}

// CutUpperCase is a split function for a [SplitFunc]
// that splits "UpperCasedWords" into "Upper" "Cased" "Words"
func CutUpperCase(s string) (tokUntil int, skipUntil int) {
	org := s
	if len(s) < 1 {
		return len(s), len(s)
	}
	s = s[1:]
	var i int
	for len(s) > 0 {
		r, j := utf8.DecodeRuneInString(s)
		i += j
		if unicode.IsUpper(r) {
			return i, i
		}
		s = s[j:]
	}
	return len(org), len(org)
}

// SplitFunc returns an iterator over sub string of s cut by [CutterFunc].
// When n > 0, SplitFunc cuts only n times and
// the returned iterator yields rest of string after n sub strings, if non empty.
// The sub strings from the iterator overlaps if splitFn decides so.
// splitFn is allowed to return negative offsets.
// In that case the returned iterator immediately yields rest of s and stops iteration.
func SplitFunc(s string, n int, splitFn CutterFunc) iter.Seq[string] {
	return func(yield func(string) bool) {
		if splitFn == nil {
			splitFn = CutNewLine
		}
		s := s
		n := n
		for len(s) > 0 {
			tokUntil, skipUntil := splitFn(s)
			if tokUntil < 0 || skipUntil < 0 {
				yield(s)
				return
			}
			if !yield(s[:tokUntil]) {
				return
			}
			s = s[skipUntil:]
			n--
			if n == 0 {
				if len(s) > 0 {
					yield(s)
				}
				return
			}
		}
	}
}
