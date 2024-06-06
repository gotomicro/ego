package xdebug

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	compName = "Test"
	addr     = "test"
	cost     = 150 * time.Millisecond
	req      = "test"
	reply    = "test"
	line     = "test"
	err      = "test"
)

func TestMakeReqResInfo(t *testing.T) {
	out := MakeReqResInfo(compName, addr, cost, req, reply)
	exp := "\x1b[32mTest\x1b[0m \x1b[32mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[34mtest\x1b[0m\n"
	assert.Equal(t, exp, out)
}

func TestMakeReqResError(t *testing.T) {
	err := "test"
	out := MakeReqResError(compName, addr, cost, req, err)
	exp := "\x1b[31mTest\x1b[0m \x1b[31mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[31mtest\x1b[0m\n"
	assert.Equal(t, exp, out)
}

func TestMakeReqResErrorV2(t *testing.T) {
	out := MakeReqResErrorV2(11, compName, addr, cost, req, "")
	exp := "\x1b[32m:0\x1b[0m \x1b[31mTest\x1b[0m \x1b[31mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[31m\x1b[0m \n"
	assert.Equal(t, exp, out)
}

func TestMakeReqAndResError(t *testing.T) {
	out := MakeReqAndResError(line, compName, addr, cost, req, err)
	exp := "\x1b[32mtest\x1b[0m \x1b[31mTest\x1b[0m \x1b[31mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[31mtest\x1b[0m"
	assert.Equal(t, exp, out)
}

func TestMakeReqAndResInfo(t *testing.T) {
	out := MakeReqAndResInfo(line, compName, addr, cost, req, reply)
	exp := "\x1b[32mtest\x1b[0m \x1b[32mTest\x1b[0m \x1b[32mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[34mtest\x1b[0m"
	assert.Equal(t, exp, out)
}
