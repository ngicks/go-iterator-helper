package hiter

import (
	"bufio"
	"iter"
)

// Scanner wraps scanner with an iterator over scanned text.
// Callers should check [bufio.Scanner.Err] after the returned iterator stops.
func Scan(scanner *bufio.Scanner) iter.Seq[string] {
	return func(yield func(text string) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}
