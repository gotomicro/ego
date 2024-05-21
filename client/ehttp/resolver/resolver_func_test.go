package resolver

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/attributes"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/server"
)

func Test_Equal(t *testing.T) {
	oldAttr := attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:    "ego",
		Scheme:  "http",
		Address: "192.168.1.1",
	})
	assert.True(t, oldAttr.Equal(attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:    "ego",
		Scheme:  "http",
		Address: "192.168.1.1",
	})))

	oldAttr2 := attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "http",
		Address:  "192.168.1.1",
		Metadata: map[string]string{"hello": "world"},
	})
	assert.True(t, oldAttr2.Equal(attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "http",
		Address:  "192.168.1.1",
		Metadata: map[string]string{"hello": "world"},
	})))

	oldAttr3 := attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "http",
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

	assert.True(t, oldAttr3.Equal(attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
		Name:     "ego",
		Scheme:   "http",
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
	})))
}

func Test_tryUpdateAttrs(t *testing.T) {
	res := &baseResolver{
		nodeInfo: make(map[string]*attributes.Attributes),
	}
	res.tryUpdateAttrs(map[string]server.ServiceInfo{
		"192.168.1.1": {
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.1",
		},
		"192.168.1.2": {
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.2",
		},
	})
	assert.Equal(t, 2, len(res.nodeInfo))
	assert.True(t, reflect.DeepEqual(res.nodeInfo, map[string]*attributes.Attributes{
		"192.168.1.1": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.1",
		}),
		"192.168.1.2": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.2",
		}),
	}))

	res.tryUpdateAttrs(map[string]server.ServiceInfo{
		"192.168.1.1": {
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.1",
		},
		"192.168.1.2": {
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.2",
		},
		"192.168.1.3": {
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.3",
		},
	})

	assert.Equal(t, 3, len(res.nodeInfo))
	assert.True(t, reflect.DeepEqual(res.nodeInfo, map[string]*attributes.Attributes{
		"192.168.1.1": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.1",
		}),
		"192.168.1.2": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.2",
		}),
		"192.168.1.3": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.3",
		}),
	}))

	res.tryUpdateAttrs(map[string]server.ServiceInfo{
		"192.168.1.1": {
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.1",
		},
		"192.168.1.3": {
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.3",
		},
	})

	assert.Equal(t, 2, len(res.nodeInfo))
	assert.True(t, reflect.DeepEqual(res.nodeInfo, map[string]*attributes.Attributes{
		"192.168.1.1": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.1",
		}),
		"192.168.1.3": attributes.New(constant.KeyServiceInfo, server.ServiceInfo{
			Name:    "svc-user",
			Scheme:  "http",
			Address: "192.168.1.3",
		}),
	}))

}
