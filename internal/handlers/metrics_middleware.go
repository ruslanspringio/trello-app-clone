package handlers

import (
	"notes-project/internal/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware(m *metrics.AppMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		m.HttpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		m.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
