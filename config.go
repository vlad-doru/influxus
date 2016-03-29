package influxus

import (
	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	PrecisionDefault        = "ns"
	DatabaseDefault         = "logrus"
	DefaultMeasurementValue = "logrus"
	BatchIntervalDefault    = 5
	BatchSizeDefault        = 100
)

type Config struct {
	Client             influx.Client
	Precision          string
	Database           string
	DefaultMeasurement string
	Tags               []string
	BatchInterval      int // seconds
	BatchSize          int
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
