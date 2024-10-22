package monitoring

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/fx"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

// NewTracerProvider sets up the tracer provider with OTLP exporter
func NewTracerProvider() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Configure OTLP exporter (modify the endpoint if necessary)
	exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	))
	if err != nil {
		return nil, err
	}

	// Set up resource information (service name, version, environment)
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("generic-integration"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.DeploymentEnvironmentKey.String("development"),
	)

	// Create the tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // In production, adjust the sampler appropriately
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)

	return tp, nil
}

func RegisterTracerShutdown(lc fx.Lifecycle, tp *sdktrace.TracerProvider) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
			return nil
		},
	})
}

// Module provides the tracing setup for Uber Fx
var Module = fx.Module("tracing",
	fx.Provide(NewTracerProvider),
	fx.Invoke(RegisterTracerShutdown),
)
