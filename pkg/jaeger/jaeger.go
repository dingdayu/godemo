package jaeger

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// 初始化变量
var (
	err    error
	Tracer opentracing.Tracer
	Closer io.Closer
)

// SpanContextKey is context span id
const SpanContextKey = "ParentSpanContext"
const TraceID = "Tracer"

// Init returns an instance of Jaeger Tracer jaeger.SamplerTypeConst,1
func Init(serviceName string, samplerType string, samplerParam float64) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  samplerType,
			Param: samplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	Tracer, Closer, err = cfg.NewTracer(config.Logger(jaeger.NullLogger))
	if err != nil {
		fmt.Printf("ERROR: cannot init Jaeger: %s\n", err.Error())
	}
	opentracing.SetGlobalTracer(Tracer)
}

// StartSpanFromHeader is start span from the request header.
func StartSpanFromHeader(header *http.Header, operationName string) (span opentracing.Span, ctx context.Context) {
	spanCtx, _ := Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*header))
	span, ctx = opentracing.StartSpanFromContext(context.Background(), operationName, ext.RPCServerOption(spanCtx))
	return span, ctx
}

// InjectSpanToHeader is inject span to the http header.
func InjectSpanToHeader(span opentracing.Span, header http.Header) {
	err := Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	if err != nil {
		fmt.Printf("ERROR: cannot inject headers: %s\n", err.Error())
	}
}

// StartSpanFromContext is start span from the context.
func StartSpanFromContext(ctx context.Context, operationName string) (span opentracing.Span, childCtx context.Context) {
	span, childCtx = opentracing.StartSpanFromContext(ctx, operationName)
	return span, childCtx
}

// StartSpanFromParentSpan is start span from the parent span.
func StartSpanFromParentSpan(parentSpan opentracing.Span, operationName string) (span opentracing.Span) {
	span = opentracing.StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
	return span
}

// GetSpanFormContext is get span from the context.
func GetSpanFormContext(ctx context.Context) (span opentracing.Span) {
	span, ok := ctx.Value(SpanContextKey).(opentracing.Span)
	if !ok {
		span = opentracing.StartSpan("<GetSpanFormContext-unkown>")
		span.Finish()
	}
	return span
}

// GetContextOfRootSpan is get root span context.
func GetContextOfRootSpan(ctx context.Context) (c context.Context) {
	return opentracing.ContextWithSpan(ctx, GetSpanFormContext(ctx))
}

// CopyContextOfRootSpan is copy root span context.
func CopyContextOfRootSpan(ctx *gin.Context) (c context.Context) {
	return GetContextOfRootSpan(ctx.Copy())
}
