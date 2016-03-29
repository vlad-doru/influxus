package main

import (
	"time"

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

func main() {
	// Using the default "logrus" measurement.
	logrus.WithFields(logrus.Fields{
		"user_id": 1,
	}).Info("User clicked")
	// Using a custom measurement.
	logrus.WithFields(logrus.Fields{
		"measurement": "click", // custom measurement
		"user_id":     1,
	}).Info("User clicked")
	// Allow batches to be flushed.
	time.Sleep(5 * time.Second)
}
