package xdebug

import (
	"testing"
	"time"
)

func TestMakeReqAndResError(t *testing.T) {
	err := MakeReqAndResError("test", "test", "test", time.Now().Sub(time.Now()), "test", "test")
	t.Log(err)
}

func TestMakeReqAndResInfo(t *testing.T) {
	err := MakeReqAndResInfo("test", "test", "test", time.Now().Sub(time.Now()), "test", "test")
	t.Log(err)
}
