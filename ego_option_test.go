package ego

import (
	"testing"
)

func TestWithHang(t *testing.T) {
	type args struct {
		hang bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "hang option set true",
			args: args{
				hang: true,
			},
			want: true,
		},
		{
			name: "hang option set true",
			args: args{
				hang: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(
				WithDisableBanner(true),
				WithHang(tt.args.hang),
			)
			if app.opts.hang != tt.want {
				t.Errorf("TestWithHang() = %v, want %v", tt.args.hang, tt.want)
			}
		})
	}
}

func TestWithDisableBanner(t *testing.T) {
	type args struct {
		disableBanner bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "disable banner option set true",
			args: args{
				disableBanner: true,
			},
			want: true,
		},
		{
			name: "disable banner option set true",
			args: args{
				disableBanner: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(
				WithDisableBanner(tt.args.disableBanner),
			)
			if app.opts.disableBanner != tt.want {
				t.Errorf("TestWithDisableBanner() = %v, want %v", tt.args.disableBanner, tt.want)
			}
		})
	}
}

func TestWithDisableFlagConfig(t *testing.T) {
	type args struct {
		disableFlagConfig bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "flag config option set true",
			args: args{
				disableFlagConfig: true,
			},
			want: true,
		},
		{
			name: "flag config option set false",
			args: args{
				disableFlagConfig: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(
				WithDisableBanner(true),
				WithDisableFlagConfig(tt.args.disableFlagConfig),
			)
			if app.opts.disableFlagConfig != tt.want {
				t.Errorf("WithDisableFlagConfig() = %v, want %v", tt.args.disableFlagConfig, tt.want)
			}
		})
	}
}
