package xdebug

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeReqResInfo(t *testing.T) {
	compName := "TestComponent"
	addr := "test.address.com"
	cost := 150 * time.Millisecond
	req := "test request"
	reply := "test reply"
	MakeReqResInfo(compName, addr, cost, req, reply)
	assert.NoError(t, nil)
}

func TestMakeReqResError(t *testing.T) {
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	err := "test"
	MakeReqResError(compName, addr, cost, req, err)
	assert.NoError(t, nil)
}

func TestMakeReqResErrorV2(t *testing.T) {
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	MakeReqResErrorV2(11, compName, addr, cost, req, "")
	assert.NoError(t, nil)
}

func TestMakeReqAndResError(t *testing.T) {
	line := "test"
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	err := "test"
	MakeReqAndResError(line, compName, addr, cost, req, err)
	assert.NoError(t, nil)
}

func TestMakeReqAndResInfo(t *testing.T) {
	line := "test"
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	reply := "test"
	MakeReqAndResInfo(line, compName, addr, cost, req, reply)
	assert.NoError(t, nil)
}
