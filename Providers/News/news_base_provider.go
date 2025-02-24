package News

import (
	"Providers/Configs"
	"context"
	"sync"
)

type IProviderNews interface {
	GetNews(ctx context.Context) ([]string, error)
	GetNewsFromRegion(ctx context.Context, region string) ([]string, error)
	GetProviderName() string
	GetConfig() Configs.ProviderConfig

	SaveNews(ctx context.Context, news string) error
	SaveNewsForRegion(ctx context.Context, news string, region string) error
	SearchNews(ctx context.Context, query string) ([]string, error)
	ResponseReaderFunc(body []byte) ([]string, error)
}

var (
	registerProviders = make(map[string]IProviderNews)
	providersMutex    sync.Mutex
)

func RegisterProvider(providerType string, provider IProviderNews) {
	// register the provider
	providersMutex.Lock()
	defer providersMutex.Unlock()
	registerProviders[providerType] = provider
}
