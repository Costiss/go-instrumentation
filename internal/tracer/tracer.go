package tracer

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	otelEndpoint := os.Getenv("OTEL_ENDPOINT")
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(otelEndpoint))
}

func newTraceProvicer(exporter sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("GoAPI"),
		),
	)
	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)
}

func InitTracer() func(context.Context) error {
	ctx := context.Background()
	exporter, err := newExporter(ctx)
	if err != nil {
		panic(err)
	}
	tracer := newTraceProvicer(exporter)
	otel.SetTracerProvider(tracer)

	return tracer.Shutdown
}
