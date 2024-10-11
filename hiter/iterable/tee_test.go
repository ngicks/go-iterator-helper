package iterable

import (
	"slices"
	"sync"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestTeePipe(t *testing.T) {
	p, seq := TeeSeqPipe(hiter.Range(0, 10))

	var (
		wg    sync.WaitGroup
		piped []int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		piped = slices.Collect(p.IntoIter())
	}()

	collected := slices.Collect(seq.IntoIter())
	wg.Wait()

	expected := slices.Collect(hiter.Range(0, 10))
	assert.DeepEqual(t, expected, collected)
	assert.DeepEqual(t, expected, piped)

	assert.Assert(t, !p.Push(5))
}
