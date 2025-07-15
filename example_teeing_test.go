package goiteratorhelper_test

import (
	"fmt"
	"sync"

	"github.com/ngicks/go-iterator-helper/hiter"
	"github.com/ngicks/go-iterator-helper/hiter/tee"
)

func Example_teeing() {
	src := hiter.Range(0, 5)
	seqPiped, seq := tee.TeeSeqPipe(0, src)

	var found bool

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		found = hiter.Contains(3, seqPiped.IntoIter())
		// Don't forget to discard all elements from seq!
		// Without this, tee could not proceed.
		hiter.Discard(seqPiped.IntoIter())
	}()

	for i := range hiter.Map(func(i int) int { return i * i }, seq.IntoIter()) {
		fmt.Printf("i = %02d\n", i)
	}
	wg.Wait()
	fmt.Printf("\nfound=%t\n", found)
	// Output:
	// i = 00
	// i = 01
	// i = 04
	// i = 09
	// i = 16
	//
	// found=true
}
