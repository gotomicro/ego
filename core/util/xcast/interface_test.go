package xcast

import (
	"testing"
)

func TestToBool(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				i: "bool",
			},
			want: false,
		},
		{
			name: "2",
			args: args{
				i: "True",
			},
			want: true,
		},
		{
			name: "3",
			args: args{
				i: "true",
			},
			want: true,
		},
		{
			name: "4",
			args: args{
				i: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBool(tt.args.i); got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
