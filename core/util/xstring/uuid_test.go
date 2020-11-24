package xstring

import (
	"testing"
	"time"
)

func TestGenerateUUID(t *testing.T) {
	type args struct {
		seedTime time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateUUID(tt.args.seedTime); got != tt.want {
				t.Errorf("GenerateUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateID(); got != tt.want {
				t.Errorf("GenerateID() = %v, want %v", got, tt.want)
			}
		})
	}
}
