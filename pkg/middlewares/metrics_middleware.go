package middlewares

import (
	"expvar"
	"time"

	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	var (
		totalRequestsReceived           = expvar.NewInt("total_requests_received")
		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
	)
	return func(c *gin.Context) {
		start := time.Now()
		totalRequestsReceived.Add(1)
		c.Next()
		totalResponsesSent.Add(1)
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	}
}
