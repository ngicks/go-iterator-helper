package ioiter_test

import (
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
