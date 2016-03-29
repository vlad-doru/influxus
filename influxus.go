package influxus

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	LevelTag     = "level"
	MessageField = "message"
	MetricField  = "metric"

	DefaultMetricField = "logrus"
)

type Hook struct {
	config Config

	// Channel used for batching.
	comm chan *influx.Point
}

func NewInfluxusHook(config *Config) (*Hook, error) {
	if config == nil {
		return nil, fmt.Errorf("Influxus configuration passed to InfluxDB is nil.")
	}
	if config.Client == nil {
		return nil, fmt.Errorf("InfluxDB client passed to Influxus configuration is nil.")
	}
	config.setDefaults()
	hook := &Hook{
		config: *config,
	}
	// Make a buffered channel so that senders will not block.
	hook.comm = make(chan *influx.Point, config.BatchSize)
	// Spawn a worker goroutine to handle the first batch.
	hook.spawnBatchHandler()
	return hook, nil
}

func (hook *Hook) spawnBatchHandler() {
	// Use a channel to control the batch interval.
	interval := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(hook.config.BatchInterval) * time.Second)
		close(interval)
	}()
	// Create a new batch locally in the handler.
	batch, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  hook.config.Database,
		Precision: hook.config.Precision,
	})
	if err != nil {
		logrus.Fatalf("Could not generate an InfluxDB batch of points: %v", err)
	}
	batchSize := 0
	for true {
		select {
		case <-interval:
			// It means we have reached the batch interval duration.
			break
		case point := <-hook.comm:
			batch.AddPoint(point)
			batchSize++
			if batchSize == hook.config.BatchSize {
				break
			}
		}
	}
	err = hook.config.Client.Write(batch)
	if err != nil {
		logrus.Errorf("Could not write batch of points to InfluxDB: %v", err)
	}
	// After we tried to write the batch we spawn a new batch handler.
	hook.spawnBatchHandler()
}

func (hook *Hook) Fire(entry *logrus.Entry) (err error) {
	// Create a new InfluxDB points and send it for processing.
	entry.Data[MessageField] = entry.Message

	metric := DefaultMetricField
	if result, ok := getTag(entry.Data, MetricField); ok {
		metric = result
	}

	tags := map[string]string{
		LevelTag: entry.Level.String(),
	}
	// Complete with the tags necessary.
	for _, tag := range hook.config.Tags {
		if tagValue, ok := getTag(entry.Data, tag); ok {
			tags[tag] = tagValue
		}
	}

	pt, err := influx.NewPoint(metric, tags, entry.Data, entry.Time)
	if err != nil {
		return fmt.Errorf("Cannot add InfluxDB point in Influxus Hook: %v", err)
	}
	// Send the point for processing.
	hook.comm <- pt
	return nil
}

// Levels implementation allows for level logging.
func (hook *Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

// Helper function.
func getTag(fields logrus.Fields, tag string) (string, bool) {
	value, ok := fields[tag]
	if ok {
		return fmt.Sprintf("%v", value), ok
	}
	return "", ok
}
