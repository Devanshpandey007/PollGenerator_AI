package news_api

import (
	Adapters "Providers/Adapters/News"
	"Providers/Configs"
	"Providers/Providers/News"
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
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
func (p *NewsApiProvider) SearchNews(ctx context.Context, query string) ([]string, error) {
	// We should use the newsAdapter to get the news
	news, err := p.Adapter.SearchNews(ctx, *p.Config, query, p.ResponseReaderFunc)
	if err != nil {
		return nil, err
	}
	return news, nil

}

// ResponseReaderFunc reads the response body and returns the news
func (p *NewsApiProvider) ResponseReaderFunc(body []byte) ([]string, error) {
	// Extract news titles
	var news []string
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		news = append(news, s.Text())
	})

	return news, nil
}
