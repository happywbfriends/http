package http_clt

import "github.com/prometheus/client_golang/prometheus"

func newCounter(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Counter {

	if len(labelsOpt) == 0 { // чтобы не засорять счетчик ссылками на пустые мапки
		labelsOpt = nil
	}

	s := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		})
	prometheus.MustRegister(s)
	return s
}

func newGauge(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Gauge {
	if len(labelsOpt) == 0 { // чтобы не засорять счетчик ссылками на пустые мапки
		labelsOpt = nil
	}

	g := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		})
	prometheus.MustRegister(g)
	return g
}

func newSummary(ns, subsystem, name string, labelsOpt map[string]string) prometheus.Summary {
	if len(labelsOpt) == 0 { // чтобы не засорять счетчик ссылками на пустые мапки
		labelsOpt = nil
	}

	s := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:   ns,
			Subsystem:   subsystem,
			Name:        name,
			ConstLabels: labelsOpt,
		})
	prometheus.MustRegister(s)
	return s
}
