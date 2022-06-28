package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	scrapeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricNameSpace,
		Name:      "gathering_duration",
		Help:      "The duration the gathering of all data took",
	}, []string{})
	hostHealth = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNameSpace,
		Name:      "host_health",
		Help:      "Indicates if a host is healthy",
	}, []string{
		"host",
		"addr",
	})
)
