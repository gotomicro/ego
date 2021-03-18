package emetric

import "github.com/prometheus/client_golang/prometheus"

// SummaryVecOpts ...
type SummaryVecOpts struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	Objectives map[float64]float64
	Labels     []string
}

// SummaryVec ...
type SummaryVec struct {
	*prometheus.SummaryVec
}

// Build ...
func (opts SummaryVecOpts) Build() *SummaryVec {
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  opts.Namespace,
			Subsystem:  opts.Subsystem,
			Name:       opts.Name,
			Help:       opts.Help,
			Objectives: opts.Objectives,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &SummaryVec{
		SummaryVec: vec,
	}
}

// Observe ...
func (summary *SummaryVec) Observe(v float64, labels ...string) {
	summary.WithLabelValues(labels...).Observe(v)
}
