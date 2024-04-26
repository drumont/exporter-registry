package exporter_registry

import "github.com/prometheus/client_golang/prometheus"

var _ prometheus.Collector = &basicCollector{}

type basicCollector struct {
	ConnectionsActive *prometheus.Desc

	stats func() ([]NginxStats, error)
}

func NewBasicCollector(stats func() ([]NginxStats, error)) prometheus.Collector {
	return &basicCollector{
		ConnectionsActive: prometheus.NewDesc(
			"nginx_connections_active",
			"Number ...",
			[]string{},
			nil,
		),
		stats: stats,
	}
}

func (c *basicCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.ConnectionsActive,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *basicCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.stats()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.ConnectionsActive, err)
		return
	}

	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			c.ConnectionsActive,
			prometheus.GaugeValue,
			s.ConnectionsActive,
		)
	}
}
