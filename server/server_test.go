package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/util/xtime"
)

func TestServiceInfo(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want ServiceInfo
	}{
		{
			name: "test",
			args: args{options: []Option{
				WithMetaData("key", "val"),
				WithScheme("http"),
				WithAddress("localhost"),
				WithName("myserver"),
				WithKind(constant.ServiceProvider),
			}},
			want: ServiceInfo{
				Name:    "myserver",
				Scheme:  "http",
				Address: "localhost",
				Weight:  100,
				Enable:  true,
				Healthy: true,
				Kind:    constant.ServiceProvider,
				Metadata: map[string]string{
					"appHost":    "",
					"appMode":    "",
					"appVersion": "",
					"buildTime":  "",
					"egoVersion": "unknown version",
					"key":        "val",
					"startTime":  xtime.TS.Format(time.Now()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyOptions(tt.args.options...)
			label := got.Label()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, "http://localhost", label)
		})
	}
}
