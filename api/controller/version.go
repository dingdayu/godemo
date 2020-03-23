package controller

import (
	"net/http"
	"time"

	"demo/pkg/enum"
	"demo/pkg/jaeger"

	"github.com/gin-gonic/gin"
)

var BuildTime string
var BuildVersion string

// @Summary 获取接口版本
// @Produce  json
// @SuccessResponse 200 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/v1/version [get]
func Version(c *gin.Context) {
	ctx := jaeger.GetContextOfRootSpan(c)
	span, ctx := jaeger.StartSpanFromContext(ctx, "Validate")
	defer span.Finish()

	if len(BuildTime) == 0 {
		BuildTime = time.Now().Format(enum.DateTimeFormat)
		BuildVersion = time.Now().Format("20060102150405")
	}

	span.LogKV("time", BuildTime)
	span.LogKV("version", BuildVersion)

	res := map[string]interface{}{}
	res["code"] = 200
	res["message"] = "success"
	res["data"] = map[string]string{
		"time":     BuildTime,
		"version":  BuildVersion,
		"trace_id": c.GetString(jaeger.TraceID),
	}
	c.JSON(http.StatusOK, res)
}
