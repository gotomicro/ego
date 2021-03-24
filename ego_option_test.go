package ego

import (
	"testing"
)

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
			name: "",
			args: args{
				disableFlagConfig: true,
			},
			want: true,
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
