package router

import (
	"demo/api/controller"
	"demo/api/controller/health"
	"demo/api/middleware"
	"demo/pkg/config"
	"demo/pkg/log"

	"github.com/gin-contrib/expvar"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var handle *gin.Engine

func Handler() *gin.Engine {
	handle = gin.New()
	handle.ForwardedByClientIP = true

	// 读取配置文件
	if config.GetString("app.model") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if gin.Mode() != gin.ReleaseMode {
		// 非正式环境下开始开启 debug 模式
		handle.GET("/debug/vars", expvar.Handler())
	}

	// 开启 Recover
	handle.Use(middleware.RecoveryWithZap(log.New().Named(middleware.RecoveryLogNamed), true))
	// 开启 gzip
	handle.Use(gzip.Gzip(gzip.DefaultCompression))

	// 探活与采集接口
	handle.GET("/", health.Hello)
	handle.HEAD("/health", health.Hello)
	handle.GET("/health", health.Hello)
	handle.GET("/ping", health.Ping)
	handle.GET("/metrics", health.Prometheus)

	v1 := handle.Group("/v1")

	// 根据配置决定是否启用 api 请求日志
	if config.GetBool("log.request_log") {
		// 启用链路追踪
		v1.Use(middleware.Jaeger())
		// 启用请求日志中间件
		v1.Use(middleware.WriterLog(log.New().Named(middleware.RequestLogNamed)))
	}

	// 服务路由
	v1.GET("/version", controller.Version)

	return handle
}
