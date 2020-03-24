package datadog

import (
	"context"
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type TelemetryOption func(*telemetry)

func WithRate(rate float64) TelemetryOption {
	return func(t *telemetry) {
		t.rate = rate
	}
}

func WithTags(tags map[string]string) TelemetryOption {
	return func(t *telemetry) {
		for k, v := range tags {
			t.tags = append(t.tags, fmt.Sprintf("%s:%s", k, v))
		}
	}
}

func WithTag(key, value string) TelemetryOption {
	return func(t *telemetry) {
		t.tags = append(t.tags, fmt.Sprintf("%s:%s", key, value))
	}
}

func WithTagsList(tags []string) TelemetryOption {
	return func(t *telemetry) {
		t.tags = tags
	}
}

type telemetry struct {
	rate float64
	tags []string
}

// Gauge measures the value of a metric at a particular time.
func Gauge(ctx context.Context, name string, value float64, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Gauge(name, value, g.tags, g.rate)
}

// Count tracks how many times something happened per second.
func Count(ctx context.Context, name string, value int64, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Count(name, value, g.tags, g.rate)
}

// Histogram tracks the statistical distribution of a set of values on each host.
func Histogram(ctx context.Context, name string, value float64, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Histogram(name, value, g.tags, g.rate)
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func Distribution(ctx context.Context, name string, value float64, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Distribution(name, value, g.tags, g.rate)
}

// Decr is just Count of -1
func Decr(ctx context.Context, name string, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Decr(name, g.tags, g.rate)
}

// Incr is just Count of 1
func Incr(ctx context.Context, name string, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Incr(name, g.tags, g.rate)
}

// Set counts the number of unique elements in a group.
func Set(ctx context.Context, name string, value string, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Set(name, value, g.tags, g.rate)
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(ctx context.Context, name string, value time.Duration, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).Timing(name, value, g.tags, g.rate)
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func TimingDefer(ctx context.Context, name string, opts ...TelemetryOption) func() {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	start := time.Now()
	return func() {
		ClientFromContext(ctx).Timing(name, time.Now().Sub(start), g.tags, g.rate)
	}
}

// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
func TimeInMilliseconds(ctx context.Context, name string, value float64, opts ...TelemetryOption) error {
	g := &telemetry{rate: 1}
	for _, opt := range opts {
		opt(g)
	}

	return ClientFromContext(ctx).TimeInMilliseconds(name, value, g.tags, g.rate)
}

// Event sends the provided Event.
func Event(ctx context.Context, e *statsd.Event) error {
	return ClientFromContext(ctx).Event(e)
}

// SimpleEvent sends an event with the provided title and text.
func SimpleEvent(ctx context.Context, title, text string) error {
	return ClientFromContext(ctx).SimpleEvent(title, text)
}

// Close the client connection.
func Close(ctx context.Context) error {
	return ClientFromContext(ctx).Close()
}

// Flush forces a flush of all the queued dogstatsd payloads.
func Flush(ctx context.Context) error {
	return ClientFromContext(ctx).Flush()
}

// SetWriteTimeout allows the user to set a custom write timeout.
func SetWriteTimeout(ctx context.Context, d time.Duration) error {
	return ClientFromContext(ctx).SetWriteTimeout(d)
}
