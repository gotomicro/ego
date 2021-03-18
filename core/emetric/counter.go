package emetric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// CounterVecOpts ...
type CounterVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

// Build ...
func (opts CounterVecOpts) Build() *CounterVec {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &CounterVec{
		CounterVec: vec,
	}
}

// NewCounterVec ...
func NewCounterVec(name string, labels []string) *CounterVec {
	return CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}

// CounterVec ...
type CounterVec struct {
	*prometheus.CounterVec
}

// Inc ...
func (counter *CounterVec) Inc(labels ...string) {
	counter.WithLabelValues(labels...).Inc()
}

// Add ...
func (counter *CounterVec) Add(v float64, labels ...string) {
	counter.WithLabelValues(labels...).Add(v)
}
