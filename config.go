package influxus

import (
	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	// PrecisionDefault represents the default precision used for the InfluxDB points.
	PrecisionDefault = "ns"
	// DatabaseDefault is the default database that we will write to, if not specified otherwise in the Config for the hook.
	DatabaseDefault = "logrus"
	// DefaultMeasurementValue is the default measurement that we will assign to each point, unless there is a field called "measurement".
	DefaultMeasurementValue = "logrus"
	// BatchIntervalDefault represents the number of seconds that we wait for a batch to fill up.
	// After that we flush it to InfluxDB whatsoever.
	BatchIntervalDefault = 5
	// BatchSizeDefault represents the maximum size of a batch.
	BatchSizeDefault = 100
)

// Config is the struct that we will use to configure our Influxus hook to Logrus.
type Config struct {
	Client             influx.Client
	Precision          string
	Database           string
	DefaultMeasurement string
	// Tags that we will extract from the log fields and set them as Influx point tags.
	Tags          []string
	BatchInterval int // seconds
	BatchSize     int
}

func (config *Config) setDefaults() {
	if config.Precision == "" {
		config.Precision = PrecisionDefault
	}
	if config.Database == "" {
		config.Database = DatabaseDefault
	}
	if config.DefaultMeasurement == "" {
		config.DefaultMeasurement = DefaultMeasurementValue
	}
	if config.BatchInterval <= 0 {
		config.BatchInterval = BatchIntervalDefault
	}
	if config.BatchSize <= 0 {
		config.BatchSize = BatchSizeDefault
	}
}
