package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegenHiddenFile(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{
			"/a/b/c/d.txt",
			"/a/b/c/.d.txt",
		},
		{
			"./a/b/c/d.txt",
			"a/b/c/.d.txt",
		},
		{
			"d.txt",
			".d.txt",
		},
	}
	for _, c := range cases {
		out := genHiddenFile(c.in)
		assert.Equal(t, c.want, out)
	}
}
