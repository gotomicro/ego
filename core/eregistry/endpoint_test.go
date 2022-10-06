package eregistry

import (
	"reflect"
	"testing"

	"github.com/gotomicro/ego/server"
	"github.com/stretchr/testify/assert"
)

func TestEndpoints_DeepCopy(t *testing.T) {
	in := newEndpoints()
	in.DeepCopy()
	assert.True(t, reflect.DeepEqual(in, in.DeepCopy()))
	var in2 *Endpoints
	assert.Nil(t, in2.DeepCopy())

}

func TestEndpoints_DeepCopyInfo(t *testing.T) {
	in := newEndpoints()
	out := newEndpoints()
	in.deepCopyInfo(out)
	assert.True(t, reflect.DeepEqual(in, out))
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
