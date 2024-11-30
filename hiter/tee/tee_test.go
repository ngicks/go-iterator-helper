package tee

import (
	"slices"
	"sync"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter"
	"gotest.tools/v3/assert"
)

func TestTeeSeqPipe(t *testing.T) {
	input := hiter.Range(0, 10)
	expected := slices.Collect(input)
	t.Run("unbuffered", func(t *testing.T) {
		p, seq := TeeSeqPipe(0, hiter.Range(0, 10))
		defer p.Close()
		defer seq.Stop()

		var (
			wg     sync.WaitGroup
			piped2 []int
		)

		wg.Add(1)
		go func() {
			defer wg.Done()
			piped2 = slices.Collect(p.IntoIter())
		}()

		piped1 := slices.Collect(seq.IntoIter())
		wg.Wait()
		assert.DeepEqual(t, expected, piped1)
		assert.DeepEqual(t, expected, piped2)

		assert.Assert(t, !p.Push(5))
	})
	t.Run("buffered", func(t *testing.T) {
		p, seq := TeeSeqPipe(1, input)
		defer p.Close()
		defer seq.Stop()

		var piped1, piped2 []int
		for {
			v1, ok2 := hiter.First(seq.IntoIter())
			v2, ok1 := hiter.First(p.IntoIter())
			if !ok1 || !ok2 {
				break
			}
			piped1 = append(piped1, v2)
			piped2 = append(piped2, v1)
		}

		assert.DeepEqual(t, expected, piped1)
		assert.DeepEqual(t, expected, piped2)
		assert.Assert(t, !p.Push(5))
	})
}
