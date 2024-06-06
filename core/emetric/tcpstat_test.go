package emetric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePort(t *testing.T) {
	//got, err := parsePort("18EB")
	//assert.Equal(t, nil, err)
	//assert.Equal(t, 6379, got)
	got1, err1 := parseIpV4("95141EAC:18EB")
	assert.Equal(t, nil, err1)
	assert.Equal(t, "172.30.20.149:6379", string(got1))
}
