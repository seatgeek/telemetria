package telemetria

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

// Recorder represents the a telegraf or influx db client that is
// configured to write to a specific endpoint one or multiple metrics
type Recorder interface {
	// Send a single metric to the recorder enpoint
	WriteOne(Metric) error

	// Send a batch of metrics to the recorder endpoint
	WriteMany([]Metric) error
}

// SimpleRecorder The simplest implementation of the Recorder interface
// It uses all the defaults from the underlying influxdb client and only
// allows configuring the precision.
type SimpleRecorder struct {
	Client    *influxdb.Client
	Database  string
	Precision string
}

// NoRecorder Represents the idea of a mocked Recorder, any call to its
// interface methods will be a noop
type NoRecorder struct{}

// Metric is a recordable stats poin identified by a name and a series of
// fields with values. Metrics can also be decorated with tags, which are
// key-value pairs that can be used for searching and aggregating related
// metrics.
type Metric struct {
	// The metric name. In influxdb this will be the table table
	Name string

	// Fields are made up of field keys and field values. Field values are your
	// data; they can be strings, floats, integers, or booleans. In influxdb
	// fields correspond to the columns in the table.
	Fields map[string]interface{}

	// The key-value pair of tags to assign to the metric. Tags are optional
	// metadata, but they are a good idea to have, since in influxdb tags are
	// indexed, so you can quickly search by them.
	Tags map[string]string
}

// NewRecorder creates a new Recorder, capable of transmiting the metrics
// to a persistent storage. The persisten storage will be located at the given
// address.
func NewRecorder(address string) (Recorder, error) {
	a, err := url.Parse(address)

	if err != nil {
		return nil, fmt.Errorf("Could not parse address '%s':\n %s", address, err)
	}

	if strings.IndexAny(a.Scheme, "http") == 0 {
		return newHTTPClient(a)
	}

	if strings.IndexAny(a.Scheme, "udp") == 0 {
		return newUDPClient(a)
	}

	return nil, fmt.Errorf("I don't know how to create a client for '%s'", address)
}

func newUDPClient(address *url.URL) (Recorder, error) {
	influxClient, err := influxdb.NewUDPClient(influxdb.UDPConfig{
		Addr: address.Host,
	})

	if err != nil {
		return nil,
			fmt.Errorf("Could not create a UDP recorder:\n %s", err.Error())
	}

	return SimpleRecorder{
		Client:    &influxClient,
		Database:  strings.Replace(address.Path, "/", "", 1),
		Precision: "ns",
	}, nil
}

func newHTTPClient(address *url.URL) (Recorder, error) {
	var user string
	var password string

	if address.User != nil {
		user = address.User.Username()
		password, _ = address.User.Password()
	}

	influxClient, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     fmt.Sprintf("%s://%s", address.Scheme, address.Host),
		Username: user,
		Password: password,
	})

	if err != nil {
		return nil,
			fmt.Errorf("Could not create a HTTP recorder:\n %s", err.Error())
	}

	return SimpleRecorder{
		Client:    &influxClient,
		Database:  strings.Replace(address.Path, "/", "", 1),
		Precision: "ns",
	}, nil
}

// WriteOne Immediately store the metric in the persistent storage
func (r SimpleRecorder) WriteOne(metric Metric) error {
	return r.WriteMany([]Metric{metric})
}

// WriteMany Immediately store the metrics in the persistent storage
func (r SimpleRecorder) WriteMany(metrics []Metric) error {
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Precision: r.Precision,
		Database:  r.Database,
	})

	if err != nil {
		return fmt.Errorf("Error creating the metric:\n %s", err.Error())
	}

	for _, m := range metrics {
		point, err := influxdb.NewPoint(m.Name, m.Tags, m.Fields, time.Now())

		if err != nil {
			return fmt.Errorf("Could not persist a metric:\n %s", err.Error())
		}

		bp.AddPoint(point)
	}

	client := *r.Client
	return client.Write(bp)
}

// WriteOne This does nothing
func (r NoRecorder) WriteOne(m Metric) error {
	return nil
}

// WriteMany This does nothing
func (r NoRecorder) WriteMany(m []Metric) error {
	return nil
}

// WithPrecision Creates a new SimpleRecorder with the specified precision
func (r SimpleRecorder) WithPrecision(precision string) Recorder {
	r.Precision = precision

	return r
}

// WithPrecision This does nothing. Returns the same NoRecorder
func (r NoRecorder) WithPrecision(precision string) Recorder {
	return r
}
