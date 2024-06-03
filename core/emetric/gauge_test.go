package emetric

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewGaugeVec(t *testing.T) {
	name := "test_"
	labels := []string{"hello_", "world_"}

	mockGaugeVec := &prometheus.GaugeVec{}
	vec := func(opts prometheus.GaugeOpts, labels []string) *prometheus.GaugeVec {
		opts = prometheus.GaugeOpts{
			Namespace:   DefaultNamespace,
			Subsystem:   "test_",
			Name:        "test_",
			Help:        "test_",
			ConstLabels: map[string]string{"hello": "world"},
		}
		assert.Equal(t, DefaultNamespace, opts.Namespace)
		assert.Equal(t, name, opts.Name)
		assert.Equal(t, name, opts.Help)
		return mockGaugeVec
	}
	out := NewGaugeVec(name, labels)
	reflect.DeepEqual(vec, out.GaugeVec)
}
