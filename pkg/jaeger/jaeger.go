package jaeger

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

// Init returns an instance of Jaeger Tracer
func Init(serviceName string) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
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

// GenerateRequestID 生成请求 ID
func GenerateRequestID() string {
	unixnaono := time.Now().UnixNano()
	randInt := rand.New(rand.NewSource(unixnaono)).Intn(9999)
	requestid := strconv.FormatInt(unixnaono, 10) + strconv.Itoa(randInt)
	return requestid
}
