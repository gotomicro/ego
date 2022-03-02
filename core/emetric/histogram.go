package emetric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// HistogramVecOpts ...
type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

// HistogramVec ...
type HistogramVec struct {
	*prometheus.HistogramVec
}

// Build ...
func (opts HistogramVecOpts) Build() *HistogramVec {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
			Buckets:   opts.Buckets,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &HistogramVec{
		HistogramVec: vec,
	}
}

// Observe ...
func (histogram *HistogramVec) Observe(v float64, labels ...string) {
	histogram.WithLabelValues(labels...).Observe(v)
}

func (histogram *HistogramVec) ObserveWithExemplar(v float64, exemplar prometheus.Labels, labels ...string) {
	histogram.WithLabelValues(labels...).(prometheus.ExemplarObserver).ObserveWithExemplar(v, exemplar)
}
