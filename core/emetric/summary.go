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

type summaryVec struct {
	*prometheus.SummaryVec
}

// Build ...
func (opts SummaryVecOpts) Build() *summaryVec {
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  opts.Namespace,
			Subsystem:  opts.Subsystem,
			Name:       opts.Name,
			Help:       opts.Help,
			Objectives: opts.Objectives,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &summaryVec{
		SummaryVec: vec,
	}
}

// Observe ...
func (summary *summaryVec) Observe(v float64, labels ...string) {
	summary.WithLabelValues(labels...).Observe(v)
}
