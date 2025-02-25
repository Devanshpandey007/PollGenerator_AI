package Trends

import (
	"Providers/Configs"
	"context"
	"encoding/xml"
	"sync"
)

type IProviderTrends interface {
	GetTrends(ctx context.Context) ([]string, error)
	GetTrend(ctx context.Context, id string) (string, error)
	GetTrendFromRegion(ctx context.Context, region string) (string, error)
	GetProviderName() string
	GetConfig() Configs.ProviderConfig

	SaveTrend(ctx context.Context, trend string) error
	SaveTrendForRegion(ctx context.Context, trend string, region string) error
	ResponseReaderFunc(body []byte) ([]string, error)
}

// RssFeed maps the <rss> root.
type RssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title         string     `xml:"title"`
	ApproxTraffic string     `xml:"ht:approx_traffic"`
	Description   string     `xml:"description"`
	Link          string     `xml:"link"`
	PubDate       string     `xml:"pubDate"`
	Picture       string     `xml:"picture"`
	PictureSource string     `xml:"picture_source"`
	NewsItems     []NewsItem `xml:"news_item"`
}

type NewsItem struct {
	Title   string `xml:"news_item_title"`
	Snippet string `xml:"news_item_snippet"`
	URL     string `xml:"news_item_url"`
	Picture string `xml:"news_item_picture"`
	Source  string `xml:"news_item_source"`
}

var (
	registerProviders = make(map[string]IProviderTrends)
	providersMutex    sync.Mutex
)

func RegisterProvider(providerType string, provider IProviderTrends) {
	// register the provider
	providersMutex.Lock()
	defer providersMutex.Unlock()
	registerProviders[providerType] = provider
}
