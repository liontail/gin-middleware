package tracer

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/liontail/gin-middleware/konst"
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
	TracerObjectEnable  bool
}

func TracerMwWithOptions(options Options, traceHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		if options.LatencyEnable {
			t := time.Now()
			// Excute after end this middleware func
			defer func() {
				latency := time.Since(t)
				log.Println(latency)
			}()
		}

		// GET traceID from header if there is none generate one with uuid
		traceID := c.Request.Header.Get(traceHeader)
		if traceID == "" {
			traceID = options.Prefix + "_" + uuid.New().String()
			c.Request.Header.Set(traceHeader, traceID)
		}

		if options.TracerObjectEnable {
			c.Set("Tracer", Tracer{
				TraceID: traceID,
				Now:     time.Now(),
			})
		}
		c.Next()

		// after request

		// access the status we are sending
		if options.SendingStatusEnable {
			status := c.Writer.Status()
			log.Println(status)
		}
	}
}

func LoggerDefault() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s %s %s] %s %-7s %s | %s %d %s | \"%s %s  %s %s\"\n",
			param.TimeStamp.Format(time.RFC3339),
			konst.RedColor, param.Request.Header.Get("zr-request-id"), konst.ResetColor,
			konst.BlueColor, param.Method, konst.ResetColor,
			konst.GreenColor, param.StatusCode, konst.ResetColor,
			param.Path,
			param.Request.Proto,
			param.Latency,
			param.ErrorMessage,
		)
	})
}
