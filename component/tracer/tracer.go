package tracer

import (
	"context"
	"fmt"
	"github.com/ebar-go/ego/utils/structure"
	"github.com/petermattis/goid"
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Instance generate the uuid for per goroutine, use to mark user requests.
type Instance struct {
	collections *structure.ConcurrentMap[string, string]
}

// key use goroutine id to generate unique identifier
func (tracer *Instance) key() string {
	return fmt.Sprintf("g%d", goid.Get())
}

// Set sets the uuid for this goroutine
func (tracer *Instance) Set(id string) {
	tracer.collections.Set(tracer.key(), id)
}

// Get returns the uuid for this goroutine, it will generate a unique string if it doesn't exist'
func (tracer *Instance) Get() (id string) {
	id, ok := tracer.collections.Get(tracer.key())
	if ok {
		return id
	}
	id = uuid.NewV4().String()
	tracer.Set(id)
	return
}

// Release remove the uuid of this goroutine
func (tracer *Instance) Release() {
	tracer.collections.Del(tracer.key())
}

func New() *Instance {
	return &Instance{collections: structure.NewConcurrentMap[string, string]()}
}

type OpenTracer struct {
	options []trace.TracerProviderOption
}

func NewOpenTracer(name string) *OpenTracer {
	return &OpenTracer{options: []trace.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		// Always be sure to batch in production.
		// Record information about this application in an Resource.
		trace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", name),
		))}}
}

func (tracer *OpenTracer) WithOption(options ...trace.TracerProviderOption) *OpenTracer {
	tracer.options = append(tracer.options, options...)
	return tracer
}

func (tracer *OpenTracer) WithJaeger(endpoint string) *OpenTracer {
	export, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		panic(err)
	}
	return tracer.WithOption(trace.WithBatcher(export))
}

func (tracer *OpenTracer) WithSampler() *OpenTracer {
	return tracer.WithOption(trace.WithSampler(trace.NeverSample()))
}

func (tracer *OpenTracer) Init() {
	tp := trace.NewTracerProvider(tracer.options...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if span := oteltrace.SpanContextFromContext(ctx); span.HasTraceID() {
		return span.TraceID().String()
	}
	return ""
}
