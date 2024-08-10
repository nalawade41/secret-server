package trace

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func SetTracing() (context.Context, *trace.TracerProvider, error) {
	ctx := context.Background()

	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		return ctx, nil, errors.Wrap(err, "failed to create tracer provider")
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return ctx, tp, err
}
