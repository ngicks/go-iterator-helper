package hiter

import (
	"bufio"
	"iter"
)

// Scanner wraps scanner with an iterator over scanned text.
// The caller should check [bufio.Scanner.Err] after the returned iterator stops
// to see if it has been stopped for an error.
func Scan(scanner *bufio.Scanner) iter.Seq[string] {
	return func(yield func(text string) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}
