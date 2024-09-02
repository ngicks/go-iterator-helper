package ih

import (
	"iter"
	"strings"
	"unicode/utf8"
)

// CollectString reduces seq to a single string.
// sizeHint hints size of internal buffer.
// Correctly sized sizeHint may reduce allocation.
func CollectString(seq iter.Seq[string], sizeHint int) string {
	var buf strings.Builder
	buf.Grow(sizeHint)
	for s := range seq {
		buf.WriteString(s)
	}
	return buf.String()
}

// StringChunk returns an iterator over non overlapping sub strings of n bytes.
// Sub slicing may cut in mid of utf8 sequences.
func StringChunk(s string, n int) iter.Seq[string] {
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

// StringRuneChunk returns an iterator over non overlapping sub strings of n utf8 characters.
func StringRuneChunk(s string, n int) iter.Seq[string] {
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
