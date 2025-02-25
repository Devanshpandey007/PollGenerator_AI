package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Common"
	"Providers/Configs"
	"Providers/Providers/Trends"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
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

// fetchTrends fetches and parses the Google Trends RSS feed,
// returning a slice of combined titles.
//
// @param ctx context.Context
// @param config Configs.ProviderConfig
//
// @return []string, error
func (ta *TrendsAdapter) fetchTrends(ctx context.Context, config Configs.ProviderConfig) ([]string, error) {
	url := buildTrendsUrl(config)

	// 1. Make the GET request
	resp, err := ta.HttpClient.MakeGetRequest(ctx, url, map[string]string{})
	if err != nil {
		// Log error and send metric
		Common.LogError("Error fetching trends", err, nil)
		ta.MetricClient.SendMetric("FetchTrendsError", 1, "Count", map[string]string{
			"Provider": config.ProviderName,
		})
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			Common.LogError("Error closing response body", cerr, nil)
		}
	}()

	// 2. Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 3. Unmarshal the XML into our RssFeed struct
	var feed Trends.RssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feed: %v", err)
	}

	// 4. Build the slice of combined titles
	//    For each item -> [Item Title, news_item_title1, news_item_title2, ...]
	var results []string
	for _, item := range feed.Channel.Items {
		parts := []string{item.Title}
		for _, news := range item.NewsItems {
			if strings.TrimSpace(news.Title) != "" {
				parts = append(parts, news.Title)
			}
		}
		// Join everything with a comma+space
		line := strings.Join(parts, ", ")
		results = append(results, line)
	}

	// Return the slice, one combined string per <item>
	return results, nil
}

func buildTrendsUrl(config Configs.ProviderConfig) string {
	return fmt.Sprintf("%s/%s/rss?geo=%s", config.BaseURL, config.Endpoint, config.Region)
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
		allTrends = append(allTrends, res...)
	}
	if len(allErrors) > 0 {
		// Log errors and send metric
		Common.LogError("Error fetching trends", fmt.Errorf("failed to fetch trends from all providers: %v", allErrors), nil)
		ta.MetricClient.SendMetric("FetchTrendsError", 1, "Count", nil)
		return nil, fmt.Errorf("failed to fetch trends from all providers: %v", allErrors)
	}

	duration := time.Since(startTime).Seconds()
	ta.MetricClient.SendMetric("FetchTrendsDuration", duration, "Seconds", nil)
	return allTrends, nil
}
