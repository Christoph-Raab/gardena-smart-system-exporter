package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const metricNameSpace = "gardena_smart_system"

var (
	endpointHealthCheckDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricNameSpace,
		Name:      "endpoint_health_duration",
		Help:      "The duration all endpoint health checks took",
	}, []string{})
	hostHealth = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNameSpace,
		Name:      "endpoint_health",
		Help:      "Indicates if a endpoint is healthy",
	}, []string{
		"endpoint",
		"addr",
	})
	locationsTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNameSpace,
		Name:      "locations_total",
		Help:      "The number of locations",
	}, []string{
		"endpoint",
	})
)
