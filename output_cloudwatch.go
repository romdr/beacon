package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

// CloudwatchInfo contains an instance of the cloudwatch client, and data common to all put requests
type CloudwatchInfo struct {
	Client     *cloudwatch.Client
	Dimensions []cloudwatch.Dimension
	Namespace  string
}

// Persistent state
var cloudwatchInfo *CloudwatchInfo

// Send to cloudwatch
func sendToCloudwatch(hostMetrics *HostMetrics, namespace string) {
	if cloudwatchInfo == nil {
		cloudwatchInfo = initCloudwatch(hostMetrics, namespace)
	}

	// Prepare the cloudwatch metrics data
	metricData := []cloudwatch.MetricDatum{
		makeMetricDatum(cloudwatchInfo.Dimensions, "CPUPct", hostMetrics.CPUPercent, cloudwatch.StandardUnitPercent),
		makeMetricDatum(cloudwatchInfo.Dimensions, "MemPct", hostMetrics.MemPercent, cloudwatch.StandardUnitPercent),
		makeMetricDatum(cloudwatchInfo.Dimensions, "Uptime", float64(hostMetrics.Uptime), cloudwatch.StandardUnitSeconds),
	}

	// Prepare a PutMetricData request
	req := cloudwatchInfo.Client.PutMetricDataRequest(&cloudwatch.PutMetricDataInput{
		Namespace:  &cloudwatchInfo.Namespace,
		MetricData: metricData,
	})

	// Send it and print only on error
	resp, err := req.Send(context.TODO())
	if err != nil {
		log.Printf("ERROR: Failed to send metrics to cloudwatch: %s: %s\n", err, resp)
	}
}

// Init the cloudwatch info (config, client, namespace, dimensions)
func initCloudwatch(hostMetrics *HostMetrics, namespace string) *CloudwatchInfo {
	// Load the default config and shared credentiala
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("ERROR: Failed to load AWS SDK config: %s\n", err.Error())
	}

	// Instantiate the client
	client := cloudwatch.New(cfg)

	// If the namespace isn't specified in the config, default to "beacon"
	if len(namespace) == 0 {
		namespace = "beacon"
	}

	// The beacon metrics dimensions are only the HostID and Hostname
	dimensions := []cloudwatch.Dimension{
		{
			Name:  aws.String("HostID"),
			Value: aws.String(hostMetrics.HostID),
		},
		{
			Name:  aws.String("Hostname"),
			Value: aws.String(hostMetrics.Hostname),
		},
	}

	return &CloudwatchInfo{
		Client:     client,
		Dimensions: dimensions,
		Namespace:  namespace,
	}
}

// Make a single MetricDatum object
func makeMetricDatum(dimensions []cloudwatch.Dimension, metricName string, value float64, unit cloudwatch.StandardUnit) cloudwatch.MetricDatum {
	return cloudwatch.MetricDatum{
		Dimensions: dimensions,
		MetricName: &metricName,
		Value:      &value,
		Unit:       unit,
	}
}
