package hiter

import "iter"

// AppendBytes appends byte slices from seq to b, and returns the extended slice.
func AppendBytes(b []byte, seq iter.Seq[[]byte]) []byte {
	for bb := range seq {
		b = append(b, bb...)
	}
	return b
}
