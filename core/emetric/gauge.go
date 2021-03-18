package emetric

import "github.com/prometheus/client_golang/prometheus"

// GaugeVecOpts ...
type GaugeVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

// GaugeVec ...
type GaugeVec struct {
	*prometheus.GaugeVec
}

// Build ...
func (opts GaugeVecOpts) Build() *GaugeVec {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &GaugeVec{
		GaugeVec: vec,
	}
}

// NewGaugeVec ...
func NewGaugeVec(name string, labels []string) *GaugeVec {
	return GaugeVecOpts{
		Namespace: DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}

// Inc ...
func (gv *GaugeVec) Inc(labels ...string) {
	gv.WithLabelValues(labels...).Inc()
}

// Add ...
func (gv *GaugeVec) Add(v float64, labels ...string) {
	gv.WithLabelValues(labels...).Add(v)
}

// Set ...
func (gv *GaugeVec) Set(v float64, labels ...string) {
	gv.WithLabelValues(labels...).Set(v)
}
