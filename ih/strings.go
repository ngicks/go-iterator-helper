package ih

import (
	"iter"
	"strings"
	"unicode"
	"unicode/utf8"
)

// All functions defined in this file must be prefixed with Strings, since
// it should fit in strings package.

// StringsCollect reduces seq to a single string.
// sizeHint hints size of internal buffer.
// Correctly sized sizeHint may reduce allocation.
func StringsCollect(seq iter.Seq[string], sizeHint int) string {
	var buf strings.Builder
	buf.Grow(sizeHint)
	for s := range seq {
		buf.WriteString(s)
	}
	return buf.String()
}

// StringsChunk returns an iterator over non overlapping sub strings of n bytes.
// Sub slicing may cut in mid of utf8 sequences.
func StringsChunk(s string, n int) iter.Seq[string] {
	return func(yield func(string) bool) {
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

// StringsRuneChunk returns an iterator over non overlapping sub strings of n utf8 characters.
func StringsRuneChunk(s string, n int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for len(s) > 0 {
			var i int
			for ii := n - 1; ii >= 0; ii-- {
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

// StringsCutterFunc is used to cut string from head.
// s[:tokUntil] is yielded through StringsSplitFunc.
// s[tokUntil:skipUntil] will be ignored.
type StringsCutterFunc func(s string) (tokUntil, skipUntil int)

// StringsCutNewLine is used with StringsSplitFunc.
// The input strings will be splitted at "\n".
// It also skips "\r" following "\n".
func StringsCutNewLine(s string) (int, int) {
	i := strings.Index(s, "\n")
	j := i + 1
	if i >= 0 && strings.HasPrefix(s[i:], "\n\r") {
		j++
	}
	return i, j
}

// StringsCutUpperCase splits "UpperCasedWords" into "Upper" "Cased" "Words"
func StringsCutUpperCase(s string) (tokUntil int, skipUntil int) {
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

// StringsSplitFunc returns an iterator over sub string of s cut by splitFn.
// When n > 0, StringsSplitFunc cuts only n times and
// the returned iterator yields rest of string after n sub strings, if non empty.
// The sub strings from the iterator overlaps if splitFn decides so.
// splitFn is allowed to return negative offsets.
// In that case the returned iterator immediately yields rest of s and stops iteration.
func StringsSplitFunc(s string, n int, splitFn StringsCutterFunc) iter.Seq[string] {
	if splitFn == nil {
		splitFn = StringsCutNewLine
	}
	return func(yield func(string) bool) {
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
