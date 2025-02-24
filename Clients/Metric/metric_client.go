package Metric

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"time"
)

// Ensure MetricClient implements MetricClientInterface
var _ MetricClientInterface = (*MetricClient)(nil)

// MetricClient allows sending custom metrics to CloudWatch
type MetricClient struct {
	client    cloudwatchiface.CloudWatchAPI
	namespace string
}

// NewMetricClient initializes and returns a new MetricClient
func NewMetricClient(namespace string, region string) (*MetricClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &MetricClient{
		client:    cloudwatch.New(sess),
		namespace: namespace,
	}, nil
}

// SendMetric sends a single custom metric to CloudWatch
func (m *MetricClient) SendMetric(metricName string, value float64, unit string, dimensions map[string]string) error {
	metricData := &cloudwatch.MetricDatum{
		MetricName: aws.String(metricName),
		Value:      aws.Float64(value),
		Unit:       aws.String(unit),
		Timestamp:  aws.Time(time.Now()),
	}

	// Add dimensions if provided
	if dimensions != nil {
		for key, val := range dimensions {
			metricData.Dimensions = append(metricData.Dimensions, &cloudwatch.Dimension{
				Name:  aws.String(key),
				Value: aws.String(val),
			})
		}
	}

	// Prepare the PutMetricDataInput
	input := &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(m.namespace),
		MetricData: []*cloudwatch.MetricDatum{metricData},
	}

	// Send the metric data
	_, err := m.client.PutMetricData(input)
	if err != nil {
		return fmt.Errorf("failed to send metric: %v", err)
	}

	return nil
}

// SendBatchMetrics sends multiple custom metrics to CloudWatch in a single request
func (m *MetricClient) SendBatchMetrics(metrics []*cloudwatch.MetricDatum) error {
	// Prepare the PutMetricDataInput
	input := &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(m.namespace),
		MetricData: metrics,
	}

	// Send the batch of metric data
	_, err := m.client.PutMetricData(input)
	if err != nil {
		return fmt.Errorf("failed to send batch metrics: %v", err)
	}

	return nil
}
