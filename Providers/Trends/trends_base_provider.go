package Trends

import (
	"Providers/Configs"
	"context"
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
