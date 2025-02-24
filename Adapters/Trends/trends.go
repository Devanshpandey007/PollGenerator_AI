package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Configs"
	"context"
	"fmt"
	"io"
	"time"
)

type TrendsAdapter struct {
	HttpClient   Api.ApiClientInterface
	S3Client     S3.S3ClientInterface
	MetricClient Metric.MetricClientInterface
}

// NewTrendsAdapter initializes and returns a new TrendsAdapter
// This is a constructor function that creates a new TrendsAdapter
//
// @param httpClient Api.ApiClientInterface
// @param s3Client S3.S3ClientInterface
// @param metricClient Metric.MetricClientInterface
//
// @return *TrendsAdapter
func NewTrendsAdapter(httpClient Api.ApiClientInterface, metricClient Metric.MetricClientInterface, s3Client S3.S3ClientInterface) *TrendsAdapter {
	return &TrendsAdapter{
		HttpClient:   httpClient,
		S3Client:     s3Client,
		MetricClient: metricClient,
	}
}

// fetchTrends fetches trends data from the API
// This function fetches trends data from the API and returns the response
//
// @param ctx context.Context
// @param config ProviderConfigs.ProviderConfig
//
// @return []byte, error
func (ta *TrendsAdapter) fetchTrends(ctx context.Context, config Configs.ProviderConfig) ([]byte, error) {
	url := buildTrendsUrl(config)
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", config.APIKey)}

	resp, err := ta.HttpClient.MakeGetRequest(ctx, url, headers)
	if err != nil {
		// Log error and send metric
		Common.LogError("Error fetching trends", err, nil)
		ta.MetricClient.SendMetric("FetchTrendsError", 1, "Count", map[string]string{
			"Provider": config.ProviderName,
		})
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Common.LogError("Error closing response body", err, nil)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func buildTrendsUrl(config Configs.ProviderConfig) string {
	return fmt.Sprintf("%s%s%s%s%s", config.BaseURL, config.Endpoint, config.TimeRange, config.QueryParams, config.Region)
}

// GetTrends fetches trends data from the API
// This function fetches trends data from the API and returns the response
//
// @param ctx context.Context
func (ta *TrendsAdapter) GetTrends(ctx context.Context, providerConfigs []Configs.ProviderConfig, responseReaderFunc func([]byte) ([]string, error)) ([]string, error) {
	startTime := time.Now()
	var allErrors []error
	var allTrends []string

	// Fetch trends data from the API
	for _, config := range providerConfigs {
		res, err := ta.fetchTrends(ctx, config)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}

		// Process trends data
		trends, err := responseReaderFunc(res)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}
		allTrends = append(allTrends, trends...)
	}
	if len(allErrors) > 0 {
		// Log errors and send metric
		Common.LogError("Error fetching trends", fmt.Errorf("failed to fetch trends from all providers: %v", allErrors), nil)
		ta.MetricClient.SendMetric("FetchTrendsError", 1, "Count", nil)
	}

	duration := time.Since(startTime).Seconds()
	ta.MetricClient.SendMetric("FetchTrendsDuration", duration, "Seconds", nil)
	return allTrends, nil
}
