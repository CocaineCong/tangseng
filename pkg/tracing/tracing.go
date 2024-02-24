package tracing

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func InitTracerProvider(Url string, serviceName string) func(ctx context.Context) error {
	ctx := context.Background()
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(Url),
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		logs.LogrusObj.Errorf("failed to init tracer, err: %v", err)
		return nil
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(newResource(serviceName)),
	)
	otel.SetTracerProvider(tp)
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3Propagator)
	otel.SetTextMapPropagator(propagator)
	return tp.Shutdown
}

func newResource(serviceName string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)
}
