package resolver

import (
	"reflect"
	"testing"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/attributes"
)

func Test_attrEqual(t *testing.T) {
	oldAttr := attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:    "ego",
		Scheme:  "grpc",
		Address: "192.168.1.1",
	})
	assert.True(t, attrEqual(oldAttr, server.ServiceInfo{
		Name:    "ego",
		Scheme:  "grpc",
		Address: "192.168.1.1",
	}))

	oldAttr2 := attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "grpc",
		Address:  "192.168.1.1",
		Metadata: map[string]string{"hello": "world"},
	})
	assert.True(t, attrEqual(oldAttr2, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "grpc",
		Address:  "192.168.1.1",
		Metadata: map[string]string{"hello": "world"},
	}))

	oldAttr3 := attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "grpc",
		Address:  "192.168.1.1",
		Metadata: map[string]string{"hello": "world"},
		Services: map[string]*server.Service{
			"hellosvc": {
				Namespace: "default",
				Name:      "first service",
				Labels:    nil,
				Methods:   nil,
			},
		},
	})
	assert.True(t, attrEqual(oldAttr3, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "grpc",
		Address:  "192.168.1.1",
		Metadata: map[string]string{"hello": "world"},
		Services: map[string]*server.Service{
			"hellosvc": {
				Namespace: "default",
				Name:      "first service",
				Labels:    nil,
				Methods:   nil,
			},
		},
	}))
}

func Test_tryUpdateAttrs(t *testing.T) {
	res := &baseResolver{
		nodeInfo: make(map[string]*attributes.Attributes),
	}
	res.tryUpdateAttrs(map[string]server.ServiceInfo{
		"192.168.1.1": {
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.1",
		},
		"192.168.1.2": {
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.2",
		},
	})
	assert.Equal(t, 2, len(res.nodeInfo))
	assert.True(t, reflect.DeepEqual(res.nodeInfo, map[string]*attributes.Attributes{
		"192.168.1.1": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.1",
		}),
		"192.168.1.2": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.2",
		}),
	}))

	res.tryUpdateAttrs(map[string]server.ServiceInfo{
		"192.168.1.1": {
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.1",
		},
		"192.168.1.2": {
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.2",
		},
		"192.168.1.3": {
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.3",
		},
	})

	assert.Equal(t, 3, len(res.nodeInfo))
	assert.True(t, reflect.DeepEqual(res.nodeInfo, map[string]*attributes.Attributes{
		"192.168.1.1": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.1",
		}),
		"192.168.1.2": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.2",
		}),
		"192.168.1.3": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.3",
		}),
	}))

	res.tryUpdateAttrs(map[string]server.ServiceInfo{
		"192.168.1.1": {
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.1",
		},
		"192.168.1.3": {
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.3",
		},
	})

	assert.Equal(t, 2, len(res.nodeInfo))
	assert.True(t, reflect.DeepEqual(res.nodeInfo, map[string]*attributes.Attributes{
		"192.168.1.1": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.1",
		}),
		"192.168.1.3": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "grpc",
			Address: "192.168.1.3",
		}),
	}))

}
