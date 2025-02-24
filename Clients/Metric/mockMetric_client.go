package Metric

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/stretchr/testify/mock"
)

// MockMetricClient is a mock implementation of MetricClient.
type MockMetricClient struct {
	mock.Mock
	client    *MockCloudWatchAPI
	namespace string
}

// NewMockMetricClient initializes a new instance of MockMetricClient.
func NewMockMetricClient(namespace string) *MockMetricClient {
	return &MockMetricClient{
		client:    &MockCloudWatchAPI{},
		namespace: namespace,
	}
}

// SendMetric is a mocked method to simulate sending a single metric to CloudWatch.
func (m *MockMetricClient) SendMetric(metricName string, value float64, unit string, dimensions map[string]string) error {
	args := m.Called(metricName, value, unit, dimensions)
	return args.Error(0)
}

// SendBatchMetrics is a mocked method to simulate sending multiple metrics to CloudWatch.
func (m *MockMetricClient) SendBatchMetrics(metrics []*cloudwatch.MetricDatum) error {
	args := m.Called(metrics)
	return args.Error(0)
}

// MockCloudWatchAPI provides a mock for CloudWatch API interactions
type MockCloudWatchAPI struct {
	cloudwatchiface.CloudWatchAPI
	mock.Mock
}

// PutMetricData is a mocked method for sending metric data to CloudWatch.
func (m *MockCloudWatchAPI) PutMetricData(input *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	args := m.Called(input)
	return &cloudwatch.PutMetricDataOutput{}, args.Error(1)
}
