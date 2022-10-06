package ego

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestWithArguments(t *testing.T) {
	//arguments default
	app := New()
	assert.Equal(t, os.Args[1:], app.opts.arguments)

	//arguments set
	app = New(
		WithArguments([]string{"--foo", "bar"}),
	)
	assert.Equal(t, []string{"--foo", "bar"}, app.opts.arguments)
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

func TestWithConfigPrefix(t *testing.T) {
	type args struct {
		configPrefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				configPrefix: "/ego",
			},
			want: "/ego",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(WithConfigPrefix(tt.args.configPrefix))
			assert.Equal(t, tt.want, app.opts.configPrefix)
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			args: args{
				timeout: 1 * time.Second,
			},
			want: 1 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(WithStopTimeout(tt.args.timeout))
			assert.Equal(t, tt.want, app.opts.stopTimeout)
		})
	}
}

func TestWithShutdownSignal(t *testing.T) {
	type args struct {
		sig os.Signal
	}
	tests := []struct {
		name string
		args args
		want os.Signal
	}{
		{
			args: args{
				sig: syscall.SIGQUIT,
			},
			want: syscall.SIGQUIT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(WithShutdownSignal(tt.args.sig))
			assert.Equal(t, tt.want, app.opts.shutdownSignals[0])
		})
	}
}
