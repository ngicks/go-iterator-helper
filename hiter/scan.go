package hiter

import (
	"bufio"
	"iter"
)

// Scanner wraps scanner with an iterator over scanned text.
func Scan(scanner *bufio.Scanner) iter.Seq2[string, error] {
	return func(yield func(text string, err error) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text(), nil) {
				return
			}
		}
		if scanner.Err() != nil {
			_ = yield("", scanner.Err())
			return
		}
	}
}
