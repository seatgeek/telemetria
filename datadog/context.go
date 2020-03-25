package datadog

import (
	"context"

	"github.com/DataDog/datadog-go/statsd"
)

type statsdContextKey string

const statsdClient = statsdContextKey("statsd.client")

// Set the statsd client to the context
func SetClient(ctx context.Context, client statsd.ClientInterface) context.Context {
	return context.WithValue(ctx, statsdClient, client)
}

// Create new statsd client and assign to context
func New(ctx context.Context, namespace string, options ...statsd.Option) context.Context {
	return SetClient(ctx, CreateClient(namespace, options...))
}

// Create new statsd client and return it
func CreateClient(namespace string, options ...statsd.Option) statsd.ClientInterface {
	options = append(options, statsd.WithNamespace(namespace))

	client, err := statsd.New("", options...)
	if err != nil {
		panic(err)
	}

	return client
}

// Extract client from context
func ClientFromContext(ctx context.Context) statsd.ClientInterface {
	value := ctx.Value(statsdClient)
	if value == nil {
		panic("No statsd client found in context")
	}

	return value.(statsd.ClientInterface)
}
