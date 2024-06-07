package emetric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var f []uint64
var ts = NewTCPStatCollector(f)

func TestParseIpV4(t *testing.T) {
	got, err := ts.parseIpV4("95141EAC:18EB")
	assert.Equal(t, nil, err)
	// assert.Equal(t, "172.30.20.149:6379", string(got))
	assert.Equal(t, "all", string(got))
}
