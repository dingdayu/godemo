package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"demo/pkg/jaeger"
	"demo/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const RequestLogNamed = "http_request"

// 处理跨域请求,支持options访问
func WriterLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 准备相应日志
		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)
		body := bodyBuf.Bytes()
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		latency := time.Since(start)
		var fs = []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Duration("latency", latency),
			zap.Any("response", blw.body.String()),
			zap.String("agent", c.Request.UserAgent()),
		}
		// Append error field if this is an erroneous request.
		if len(c.Errors) > 0 {
			fs = append(fs, zap.Strings("errors", c.Errors.Errors()))
		}
		// Post Method writer body
		if c.Request.Method != http.MethodGet {
			fs = append(fs, zap.ByteString("body", body))
		}
		// Writer X-Request-Id to log
		xRequestId := c.Request.Header.Get("X-Request-Id")
		if len(xRequestId) > 0 {
			fs = append(fs, zap.String("request_id", xRequestId))
		}
		// Writer trace_id to log
		if traceID := c.Value(jaeger.TraceID); traceID != nil {
			fs = append(fs, zap.String(log.TraceID, traceID.(string)))
		}
		logger.Info(c.Request.RequestURI, fs...)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
