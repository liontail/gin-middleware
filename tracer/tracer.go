package tracer

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Tracer represent state for each request.
type Tracer struct {
	TraceID string
	Now     time.Time
}

func TracerMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		// t := time.Now()

		// GET traceID from header if there is none generate one with uuid
		traceID := c.Request.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("Tracer", Tracer{
			TraceID: traceID,
			Now:     time.Now(),
		})
		c.Next()
		// // after request

		// latency := time.Since(t)
		// log.Println(latency)

		// // access the status we are sending
		// status := c.Writer.Status()
		// log.Println(status)
	}
}
