package emetric

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestCounterVec(t *testing.T) {
	name := "test"
	labels := []string{"hello", "world"}

	mockCounterVec := &prometheus.CounterVec{}
	vec := func(opts prometheus.GaugeOpts, labels []string) *prometheus.CounterVec {
		opts = prometheus.GaugeOpts{
			Namespace:   DefaultNamespace,
			Subsystem:   "test",
			Name:        "test",
			Help:        "test",
			ConstLabels: map[string]string{"hello": "world"},
		}
		assert.Equal(t, "ego", opts.Namespace)
		assert.Equal(t, name, opts.Name)
		assert.Equal(t, name, opts.Help)
		return mockCounterVec
	}
	out := NewCounterVec(name, labels)
	reflect.DeepEqual(vec, out.CounterVec)
}
