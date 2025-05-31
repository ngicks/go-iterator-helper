package ioiter_test

import (
	"errors"
	"io/fs"
	"os"
	"slices"
	"testing"

	"github.com/ngicks/go-iterator-helper/hiter/ioiter"
	"gotest.tools/v3/assert"
)

func TestReaddir(t *testing.T) {
	dir, err := os.Open("./testdata/readdir")
	if err != nil {
		panic(err)
	}
	var names []string
	for dirent, err := range ioiter.Readdir(dir) {
		if err != nil {
			panic(err)
		}
		names = append(names, dirent.Name())
	}
	// maybe unordered
	slices.Sort(names)
	assert.DeepEqual(t, []string{"bar", "baz", "foo"}, names)
}

// Test early termination
func TestReaddirEarlyTermination(t *testing.T) {
	dir, err := os.Open("./testdata/readdir")
	if err != nil {
		panic(err)
	}

	var names []string
	for dirent, err := range ioiter.Readdir(dir) {
		if err != nil {
			panic(err)
		}
		names = append(names, dirent.Name())
		if len(names) >= 2 {
			break // This should trigger the early termination path
		}
	}
	assert.Assert(t, len(names) >= 2)
}

// Mock reader that returns an error
type errorReader struct{}

func (e errorReader) Readdir(n int) ([]fs.FileInfo, error) {
	return nil, errors.New("mock error")
}

// Test error handling path
func TestReaddirError(t *testing.T) {
	reader := errorReader{}

	for dirent, err := range ioiter.Readdir(reader) {
		assert.Assert(t, dirent == nil)
		assert.Error(t, err, "mock error")
		break // We expect the first iteration to have an error
	}
}
