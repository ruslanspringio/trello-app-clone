package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type AppMetrics struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
}

func NewAppMetrics() *AppMetrics {
	return &AppMetrics{
		HttpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Общее количество HTTP запросов.",
			},
			[]string{"method", "path", "status_code"},
		),
		HttpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Длительность HTTP запросов в секундах.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}
