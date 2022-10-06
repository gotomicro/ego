package xnet

import "testing"

func TestGetLocalMainIP(t *testing.T) {
	ip, port, err := GetLocalMainIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip, port)
}
