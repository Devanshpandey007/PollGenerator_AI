package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Configs"
	"context"
	"fmt"
	"io"
)

type NewsAdapter struct {
	HttpClient   Api.ApiClientInterface
	MetricClient Metric.MetricClientInterface
}

// NewNewsAdapter initializes and returns a new NewsAdapter
// This is a constructor function that creates a new NewsAdapter
//
// @param httpClient Api.ApiClientInterface
// @param metricClient Metric.MetricClientInterface
//
// @return *NewsAdapter
func NewNewsAdapter(httpClient Api.ApiClientInterface, metricClient Metric.MetricClientInterface) *NewsAdapter {
	return &NewsAdapter{
		HttpClient:   httpClient,
		MetricClient: metricClient,
	}
}

// fetchNews fetches news data from the API
// This function fetches news data from the API and returns the response
//
// @param ctx context.Context
// @param config ProviderConfigs.ProviderConfig
//
// @return []byte, error
func (na *NewsAdapter) fetchNews(ctx context.Context, config Configs.ProviderConfig) ([]byte, error) {
	url := buildNewsUrl(config)
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", config.APIKey)}

	resp, err := na.HttpClient.MakeGetRequest(ctx, url, headers)
	if err != nil {
		// Log error and send metric
		Common.LogError("Error fetching news", err, nil)
		na.MetricClient.SendMetric("FetchNewsError", 1, "Count", map[string]string{
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

func buildNewsUrl(config Configs.ProviderConfig) string {

}
