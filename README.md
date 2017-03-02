# telemetira

This package aims to make is simple to connect and send metrics to influxdb, either by using the
http or udp protocol.

The official influxdb client is good enough, yet a bit verbose when trying to quickly send stats
to the database. Telemetria offers a friendlier face for performing those operations, without exposing
the richer and more complex internal API from influxdb.

## Installation

Install it using the `go get` command:

    go get github.com/seatgeek/telemetria

## Usage

You first need to create a `Recorder` that can be used later on for sending metrics to the database.
Recorders support the http and udp protocol. In order to create a client to you need to pass a
URL string accordingly.

### Creating a recorder using HTTP


```go
import (
	"github/seatgeek/telemetria"
)

recorder := telemetria.NewRecorder("http://user:pass@localhost:8086/my_database")
```

### Creating a recorder using UDP

```go
import (
	"github/seatgeek/telemetria"
)

recorder := telemetria.NewRecorder("udp://localhost:8089")
```

Note that there is no way to specify a database when using the UDP protocol. This needs to be handled
directly in the UDP configuration for influxdb or telegraf.

### Sending metrics

Metrics are created using the `Metric` struct and passing them to `WriteOne` or `WriteMany`. Theses structs
have fields for specifying both the values and tags that should be stored in the series table.

Here's an example showing all supported properties:

```go
recorder.WriteOne(Metric{
	Name: "cpu_usage",
	Fields: map[string]interface{}{
		"idle":   10.1,
        "system": 53.3,
        "user":   46.6,
	},
	Tags: map[string]string{"cpu": "cpu-total"}
})
```

You batch-send many metrics using the `WriteMany` method:

```go
metrics := []Metric{ metric1, metric2, metric3 }
recorder.WriteMany(metrics)
```
