package news_api

import (
	Adapters "Providers/Adapters/News"
	"Providers/Configs"
	"Providers/Providers/News"
	"context"
	"encoding/json"
	"fmt"
)

// NewsApiProvider implement the base provider interface
type NewsApiProvider struct {
	News.IProviderNews
	Config  *Configs.Config
	Adapter *Adapters.NewsAdapter
}

// Ensure that NewsApiProvider implements the IProviderNews interface
var _ News.IProviderNews = (*NewsApiProvider)(nil)

// NewProvider initializes and registers a NewsApiProvider instance.
func NewProvider(config *Configs.Config, news *Adapters.NewsAdapter) *NewsApiProvider {
	provider := &NewsApiProvider{
		Config:  config,
		Adapter: news,
	}
	News.RegisterProvider("Providers.NewsApiProvider", provider)
	return provider
}

// SearchNews fetches news data from the API
func (p *NewsApiProvider) SearchNews(ctx context.Context, query string) ([]News.NewsArticle, error) {
	// We should use the newsAdapter to get the news
	news, err := p.Adapter.SearchNews(ctx, *p.Config, query, p.ResponseReaderFunc)
	if err != nil {
		return nil, err
	}
	return news, nil

}

// ResponseReaderFunc reads the response body and returns the news articles.
// It extracts the title, description, publishedAt, and content from the JSON response.
func (p *NewsApiProvider) ResponseReaderFunc(body []byte) ([]News.NewsArticle, error) {
	// Define a struct that matches the JSON structure of the API response.
	var apiResponse struct {
		Status   string             `json:"status"`
		Articles []News.NewsArticle `json:"articles"`
	}

	// Unmarshal the JSON response.
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error parsing news API response: %v", err)
	}

	return apiResponse.Articles, nil
}
