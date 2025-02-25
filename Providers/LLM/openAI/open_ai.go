package openAI

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Common"
	"Providers/Configs"
	"Providers/Providers/LLM"
	"Providers/Providers/LLM/Generator"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type OpenAIProvider struct {
	LLM.IProviderLLM
	Config       *Configs.Config
	HttpClient   Api.ApiClientInterface
	MetricClient Metric.MetricClientInterface
}

// Ensure OpenAIProvider implements IProviderLLM
var _ LLM.IProviderLLM = (*OpenAIProvider)(nil)

// NewProvider initializes and registers a OpenAIProvider instance.
func NewProvider(config *Configs.Config, httpClient Api.ApiClientInterface, metricClient Metric.MetricClientInterface) *OpenAIProvider {
	provider := &OpenAIProvider{
		Config:       config,
		HttpClient:   httpClient,
		MetricClient: metricClient,
	}
	LLM.RegisterProvider("Providers.OpenAIProvider", provider)
	return provider
}

// GetConfig returns the provider configuration
func (p *OpenAIProvider) GetConfig() Configs.ProviderConfig {
	return *p.Config.ProviderConfig
}

// GetProviderName returns the name of the provider
func (p *OpenAIProvider) GetProviderName() string {
	return "OpenAIProvider"
}

// GeneratePollQuestions generates poll questions based on the topic and context summary
func (p *OpenAIProvider) GeneratePollQuestions(ctx context.Context, topic string, contextSummary string) (string, error) {
	prompt := fmt.Sprintf(Generator.PredictionMarketPrompt, topic, contextSummary)
	response, err := p.makePromptRequest(ctx, prompt)
	if err != nil {
		Common.LogError("Error generating poll questions", err, nil)
		return "", err
	}

	return response, nil
}

func (p *OpenAIProvider) makePromptRequest(ctx context.Context, prompt string) (string, error) {
	url := buildPromptUrl(p.Config.ProviderConfig)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", p.Config.ProviderConfig.APIKey),
		"Content-Type":  "application/json",
	}
	requestBody := map[string]interface{}{
		"model": p.Config.ProviderConfig.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  p.Config.DefaultConfig.MaxTokens,
		"temperature": p.Config.DefaultConfig.Temperature,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	response, err := p.HttpClient.MakePostRequest(ctx, url, headers, jsonBody)
	if err != nil {
		// Log error and send metric
		Common.LogError("Error fetching news", err, nil)
		p.MetricClient.SendMetric("FetchPromptError", 1, "Count", map[string]string{
			"Provider": p.GetProviderName(),
		})
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		p.MetricClient.SendMetric("FetchPromptError", 1, "Count", map[string]string{
			"Provider": p.GetProviderName(),
		})
		return "", err
	}

	var openAIResp LLM.OpenAIResp
	if err := json.Unmarshal(responseBody, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal openAI response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices found in openAI response")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func buildPromptUrl(config *Configs.ProviderConfig) string {
	return fmt.Sprintf("%s/%s", config.BaseURL, config.Endpoint)

}
