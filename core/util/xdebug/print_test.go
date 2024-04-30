package xdebug

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestMakeReqResInfo(t *testing.T) {
	compName := "TestComponent"
	addr := "test.address.com"
	cost := 150 * time.Millisecond
	req := "test request"
	reply := "test reply"

	expectedOutput := "\x1b[32mTestComponent\x1b[0m \x1b[32mtest.address.com\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest request\x1b[0m => \x1b[34mtest reply\x1b[0m\n"
	actualOutput := MakeReqResInfo(compName, addr, cost, req, reply)
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch. Expected: %s, Got: %s", expectedOutput, actualOutput)
	}
}

func TestMakeReqResError(t *testing.T) {
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	reply := "test"

	expectedOutput := "\x1b[31mTest\x1b[0m \x1b[31mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[31mtest\x1b[0m\n"
	actualOutput := MakeReqResError(compName, addr, cost, req, reply)
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch. Expected: %s, Got: %s", expectedOutput, actualOutput)
	}
}

func TestMakeReqResInfoV2(t *testing.T) {
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	reply := "test"
	_, file, line, _ := runtime.Caller(1)
	caller := file + ":" + strconv.Itoa(line)

	expectedOutput := fmt.Sprintf("\u001B[32m%v\u001B[0m \x1b[32mTest\x1b[0m \x1b[32mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[34mtest\x1b[0m\n", caller)
	actualOutput := MakeReqResInfoV2(1, compName, addr, cost, req, reply)
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch.\n Expected: %s\n Got: %s", expectedOutput, actualOutput)
	}
}

func TestMakeReqResErrorV2(t *testing.T) {
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	reply := "test"

	_, file, line, _ := runtime.Caller(1)
	caller := file + ":" + strconv.Itoa(line)

	expectedOutput := fmt.Sprintf("\x1b[32m%v\x1b[0m \x1b[31mTest\x1b[0m \x1b[31mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[31mtest\x1b[0m\n", caller)
	actualOutput := MakeReqResErrorV2(1, compName, addr, cost, req, reply)
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch.\n Expected: %s\n Got: %s", expectedOutput, actualOutput)
	}
}

func TestMakeReqAndResError(t *testing.T) {
	line := "test"
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	reply := "test"

	expectedOutput := "\x1b[32mtest\x1b[0m \x1b[31mTest\x1b[0m \x1b[31mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[31mtest\x1b[0m\n"
	actualOutput := MakeReqAndResError(line, compName, addr, cost, req, reply)
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch.\n Expected: %s\n Got: %s", expectedOutput, actualOutput)
	}
}

func TestMakeReqAndResInfo(t *testing.T) {
	line := "test"
	compName := "Test"
	addr := "test"
	cost := 150 * time.Millisecond
	req := "test"
	reply := "test"

	expectedOutput := "\x1b[32mtest\x1b[0m \x1b[32mTest\x1b[0m \x1b[32mtest\x1b[0m \x1b[33m[150ms]\x1b[0m \x1b[34mtest\x1b[0m => \x1b[34mtest\x1b[0m\n"
	actualOutput := MakeReqAndResInfo(line, compName, addr, cost, req, reply)
	if actualOutput != expectedOutput {
		t.Errorf("Output mismatch.\n Expected: %s\n Got: %s", expectedOutput, actualOutput)
	}
}
