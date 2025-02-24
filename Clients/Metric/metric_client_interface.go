package Metric

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// MetricClientInterface defines the contract for sending metrics to CloudWatch
type MetricClientInterface interface {
	SendMetric(metricName string, value float64, unit string, dimensions map[string]string) error
	SendBatchMetrics(metrics []*cloudwatch.MetricDatum) error
}
