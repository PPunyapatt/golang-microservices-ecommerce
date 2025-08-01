package helper

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(serviceName string) func(context.Context) error {
	// ctx := context.Background()
	// exp, err := otlptracegrpc.New(ctx,
	// 	otlptracegrpc.WithInsecure(),
	// 	otlptracegrpc.WithEndpoint("http://localhost:14268/api/traces"),
	// )
	// if err != nil {
	// 	log.Fatalf("failed to create OTLP exporter: %v", err)
	// }

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint("http://localhost:14268/api/traces"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown
}
