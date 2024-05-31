package eregistry

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/server"
)

func TestEndpoints_DeepCopy(t *testing.T) {
	in := newEndpoints()
	if in == nil {
		return
	}
	in.DeepCopy()
	// assert.True(t, reflect.DeepEqual(in, in.DeepCopy()))
	assert.Equal(t, in, in.DeepCopy())
	var in2 *Endpoints
	assert.Nil(t, in2.DeepCopy())

}

func TestEndpoints_DeepCopyInfo(t *testing.T) {
	in := newEndpoints()
	out := newEndpoints()
	for key, info := range in.Nodes {
		out.Nodes[key] = info
	}
	for key, config := range in.RouteConfigs {
		out.RouteConfigs[key] = config
	}
	for key, config := range in.ConsumerConfigs {
		out.ConsumerConfigs[key] = config
	}
	for key, config := range in.ProviderConfigs {
		out.ProviderConfigs[key] = config
	}
	in.deepCopyInfo(out)
	// assert.True(t, reflect.DeepEqual(in, out))
	assert.Equal(t, in, out)
}

func Test_newEndpoints(t *testing.T) {
	got := newEndpoints()
	assert.True(t, reflect.DeepEqual(got, &Endpoints{
		Nodes:           make(map[string]server.ServiceInfo),
		RouteConfigs:    make(map[string]RouteConfig),
		ConsumerConfigs: make(map[string]ConsumerConfig),
		ProviderConfigs: make(map[string]ProviderConfig),
	}))
}
