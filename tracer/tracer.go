package tracer

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Tracer represent state for each request.
type Tracer struct {
	TraceID string    `json:"trace_id" bson:"trace_id"`
	Now     time.Time `json:"requested_time" bson:"requested_time"`
}

func TracerMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		// t := time.Now()

		// GET traceID from header if there is none generate one with uuid
		traceID := c.Request.Header.Get("zr-request-id")
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

type Options struct {
	Prefix              string
	LatencyEnable       bool
	SendingStatusEnable bool
}

func TracerMwWithOptions(options Options, traceHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()

		// GET traceID from header if there is none generate one with uuid
		traceID := c.Request.Header.Get(traceHeader)
		if traceID == "" {
			traceID = options.Prefix + "_" + uuid.New().String()
		}
		c.Set("Tracer", Tracer{
			TraceID: traceID,
			Now:     time.Now(),
		})
		c.Next()

		// after request

		if options.LatencyEnable {
			latency := time.Since(t)
			log.Println(latency)
		}

		// access the status we are sending
		if options.SendingStatusEnable {
			status := c.Writer.Status()
			log.Println(status)
		}
	}
}
