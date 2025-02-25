package LLM

import (
	"Providers/Configs"
	"context"
	"sync"
)

type IProviderLLM interface {
	GetConfig() Configs.ProviderConfig
	GetProviderName() string
	GeneratePollQuestions(ctx context.Context, topic, contextSummary string) (string, error)
}

var (
	registerProviders = make(map[string]IProviderLLM)
	providersMutex    sync.Mutex
)

func RegisterProvider(providerType string, provider IProviderLLM) {
	// register the provider
	providersMutex.Lock()
	defer providersMutex.Unlock()
	registerProviders[providerType] = provider
}
