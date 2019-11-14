package middleware

import (
	"demo/pkg/jaeger"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/ext"
	ujaeger "github.com/uber/jaeger-client-go"
)

// Jaeger 为每个方法开启 tracing
func Jaeger() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := jaeger.StartSpanFromHeader(&c.Request.Header, c.Request.URL.Path)

		var traceID string
		if sc, ok := span.Context().(ujaeger.SpanContext); ok {
			traceID = sc.TraceID().String()
			c.Header(jaeger.TraceID, traceID)
			c.Set(jaeger.TraceID, traceID)
		}

		xRequestId := c.Request.Header.Get("X-Request-Id")
		if len(xRequestId) > 0 {
			span.SetTag("request.id", xRequestId)
		}
		c.Set(jaeger.SpanContextKey, span)

		c.Next()

		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		span.Finish()
	}
}
