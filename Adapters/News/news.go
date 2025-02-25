package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Common"
	"Providers/Configs"
	"context"
	"fmt"
	"io"
	"time"
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
func (na *NewsAdapter) fetchNews(ctx context.Context, url string, headers map[string]string, config Configs.ProviderConfig) ([]byte, error) {
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

func buildNewsUrl(config Configs.Config, topic string) string {
	return fmt.Sprintf("%s/%s?q=%s&from=%s&sortBy=publishedAt&apiKey=%s", config.ProviderConfig.BaseURL, config.ProviderConfig.Endpoint, topic, config.DefaultConfig.SearchStartDate, config.ProviderConfig.APIKey)
}

// SearchNews fetches news data from the API for a given topic.
// This function fetches news data from the API for a given topic and returns the response.
//
// @param ctx context.Context
// @param config ProviderConfigs.ProviderConfig
// @param responseReader ResponseReaderFunc
// @param topic string
//
// @return []byte, error
func (na *NewsAdapter) SearchNews(ctx context.Context, config Configs.Config, topic string, responseReader func(res []byte) ([]string, error)) ([]string, error) {
	startTime := time.Now()
	url := buildNewsUrl(config, topic)
	resp, err := na.fetchNews(ctx, url, map[string]string{}, *config.ProviderConfig)
	if err != nil {
		return nil, err
	}

	headlines, err := responseReader(resp)
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime).Seconds()
	na.MetricClient.SendMetric("FetchTrendsDuration", duration, "Seconds", nil)
	return headlines, nil
}

//// ResponseReaderFunc is a function that reads the response and returns the data.
//func ResponseReaderFunc(res []byte) ([]string, error) {
//	// Extract news titles
//	var news []string
//	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(res))
//	doc.Find("title").Each(func(i int, s *goquery.Selection) {
//		news = append(news, s.Text())
//	})
//
//	return news, nil
//}
