Influxus [![GoDoc](https://godoc.org/github.com/vlad-doru/influxus?status.svg)](https://godoc.org/github.com/vlad-doru/influxus)
===

Golang Hook for the [Logrus](https://github.com/Sirupsen/logrus) logging library, in order to output logs to [InfluxDB](https://influxdata.com/).

Usage
---

### Installation

```
go get github.com/vlad-doru/influxus
```

### Example
```go

import (
	"github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/vlad-doru/influxus"
)

func init() {
	// Create the InfluxDB client.
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		logrus.Fatalf("Error while creating the client: %v", err)
	}
	// Create and add the hook.
	hook, err := influxus.NewHook(
		&influxus.Config{
			Client:             influxClient,
			Database:           "logrus", // DATABASE MUST BE CREATED
			DefaultMeasurement: "logrus",
			BatchSize:          1, // default is 100
			BatchInterval:      1, // default is 5 seconds
		})
	if err != nil {
		logrus.Fatalf("Error while creating the hook: %v", err)
	}
	// Add the hook to the standard logger.
	logrus.StandardLogger().Hooks.Add(hook)
}
```

Concepts
---

### Concurrency

We use a non-blocking model for the hook. Each time an entry is fired, we create an InfluxDB point and put that through a buffered channel. It is then picked up by a worker goroutine that will handle the current batch and write everything once the batch hit the size or time limit. 

### Configuration

We chose to take in as a parameter the InfluxDB client so that we have greater flexibility over the way we connect to the Influx Server.
All the defaults and configuration options can be found in the GoDoc.

### Database Handling

The database should have been previously created. This is by design, as we want to avoid people generating lots of databases.
When passing an empty string for the InfluxDB database name, we default to "logrus" as the database name.

### Message Field

We will insert your message into InfluxDB with the field message.

TODO
---

- [x] Concurrent, non-blocking design.
- [  ] Add unit tests
- [  ] Set up continous integration.
