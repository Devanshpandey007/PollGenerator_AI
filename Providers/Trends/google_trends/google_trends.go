package google_trends

import (
	Adapters "Providers/Adapters/Trends"
	"Providers/Configs"
	"Providers/Providers/Trends"
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// GoogleTrendsProvider implement the base provider interface
type GoogleTrendsProvider struct {
	Trends.IProviderTrends
	Config  *Configs.Config
	Adapter *Adapters.TrendsAdapter
}

// Ensure that GoogleTrendsProvider implements the IProviderTrends interface
var _ Trends.IProviderTrends = (*GoogleTrendsProvider)(nil)

// NewProvider initializes and registers a GoogleTrendsProvider instance.
func NewProvider(config *Configs.Config, trends *Adapters.TrendsAdapter) *GoogleTrendsProvider {
	provider := &GoogleTrendsProvider{
		Config:  config,
		Adapter: trends,
	}
	Trends.RegisterProvider("Providers.GoogleTrendsProvider", provider)
	return provider
}

// GetTrends fetches trends data from the API
func (p *GoogleTrendsProvider) GetTrends(ctx context.Context) ([]string, error) {
	// We should use the trendAdapter to get the trends
	trends, err := p.Adapter.GetTrends(ctx, []Configs.ProviderConfig{p.GetConfig()}, p.ResponseReaderFunc)
	if err != nil {
		return nil, err
	}
	p.Adapter.MetricClient.SendMetric("NumberOfTrends", float64(len(trends)), "Count", map[string]string{
		"Provider": p.GetProviderName(),
	})
	return trends, nil
}

// GetTrend fetches a single trend from the API
func (p *GoogleTrendsProvider) GetTrend(ctx context.Context, id string) (string, error) {
	return "", nil
}

// GetTrendFromRegion fetches a single trend from the API based on a region
func (p *GoogleTrendsProvider) GetTrendFromRegion(ctx context.Context, region string) (string, error) {
	return "", nil
}

// GetProviderName returns the name of the provider
func (p *GoogleTrendsProvider) GetProviderName() string {
	return "GoogleTrendsProvider"
}

// SaveTrend saves a trend to the provider
func (p *GoogleTrendsProvider) SaveTrend(ctx context.Context, trend string) error {
	return nil
}

// SaveTrendForRegion saves a trend to the provider for a specific region
func (p *GoogleTrendsProvider) SaveTrendForRegion(ctx context.Context, trend string, region string) error {
	return nil
}

// GetConfig returns the providers configuration
func (p *GoogleTrendsProvider) GetConfig() Configs.ProviderConfig {
	return *p.Config.ProviderConfig
}

// ResponseReaderFunc extracts trends from Google Trends HTML response.
func (p *GoogleTrendsProvider) ResponseReaderFunc(body []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing Google Trends response: %v", err)
	}

	var trends []string
	doc.Find("div.summary-text").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			// If the text contains commas, split it into individual trends.
			if strings.Contains(text, ",") {
				parts := strings.Split(text, ",")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if part != "" {
						trends = append(trends, part)
					}
				}
			} else {
				trends = append(trends, text)
			}
		}
	})

	if len(trends) == 0 {
		return nil, fmt.Errorf("no trends found in response")
	}

	return trends, nil
}
