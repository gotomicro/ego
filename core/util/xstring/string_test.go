package xstring

import (
	"reflect"
	"testing"
)

func TestKickEmpty(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name string
		args args
		want Strings
	}{
		// TODO: Add test cases.
		{
			name: "testing",
			args: args{
				ss: []string{"", "1", "2", ""},
			},
			want: Strings{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KickEmpty(tt.args.ss); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KickEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyBlank(t *testing.T) {
	type args struct {
		ss []string
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
				ss: []string{"", "1", "2", ""},
			},
			want: true,
		},
		{
			name: "1",
			args: args{
				ss: []string{"1", "2"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyBlank(tt.args.ss); got != tt.want {
				t.Errorf("AnyBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}
