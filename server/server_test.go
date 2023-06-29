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
					"depEnv":     "",
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

func Test_ServiceInfoLabel(t *testing.T) {
	svc := ServiceInfo{
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
			"egoVersion": "v0.7.0",
			"key":        "val",
			"startTime":  xtime.TS.Format(time.Now()),
		},
	}
	assert.Equal(t, "http://localhost", svc.Label())
}

func Test_ServiceInfoEqual(t *testing.T) {
	svc := ServiceInfo{
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
			"egoVersion": "v0.7.0",
			"key":        "val",
		},
	}
	assert.True(t, svc.Equal(ServiceInfo{
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
			"egoVersion": "v0.7.0",
			"key":        "val",
		},
	}))
}

func Test_ServiceInfoGetServiceValue(t *testing.T) {
	svc := ServiceInfo{
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
			"egoVersion": "v0.7.0",
			"key":        "val",
			"startTime":  xtime.TS.Format(time.Now()),
		},
	}
	assert.Contains(t, svc.GetServiceValue(), "v0.7.0")
	assert.Contains(t, svc.GetServiceValue(), "localhost")
	assert.Contains(t, svc.GetServiceValue(), "myserver")
	assert.Contains(t, svc.GetServiceValue(), "http")
}

func Test_ServiceInfoGetServiceKey(t *testing.T) {
	svc := ServiceInfo{
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
			"egoVersion": "v0.7.0",
			"key":        "val",
			"startTime":  xtime.TS.Format(time.Now()),
		},
	}
	assert.Equal(t, "/ego/myserver/providers/http://localhost", svc.GetServiceKey("ego"))
}
