package eregistry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNop_Close(t *testing.T) {
	n := Nop{}
	err := n.Close()
	assert.NoError(t, err)
}

func TestNop_ListServices(t *testing.T) {
	n := Nop{}
	assert.Panics(t, func() {
		_, _ = n.ListServices(context.Background(), Target{})
	})
}

func TestNop_RegisterService(t *testing.T) {
	n := Nop{}
	err := n.RegisterService(context.Background(), nil)
	assert.NoError(t, err)
}

func TestNop_SyncServices(t *testing.T) {
	n := Nop{}
	err := n.SyncServices(context.Background(), SyncServicesOptions{})
	assert.NoError(t, err)
}

func TestNop_UnregisterService(t *testing.T) {
	n := Nop{}
	err := n.UnregisterService(context.Background(), nil)
	assert.NoError(t, err)
}

func TestNop_WatchServices(t *testing.T) {
	n := Nop{}
	assert.Panics(t, func() {
		_, _ = n.WatchServices(context.Background(), Target{})
	})
}
