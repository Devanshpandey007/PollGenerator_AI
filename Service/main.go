package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	ta "Providers/Adapters/Trends"
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Common"
	"Providers/Configs"
	"Providers/Providers/LLM/openAI"
	"Providers/Providers/Trends/google_trends"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/spf13/viper"
)

// LoadConfig loads configuration via Viper and returns the Viper instance
// along with an app-wide configuration struct.
func LoadConfig() (*viper.Viper, *Configs.Config, error) {
	environment, found := os.LookupEnv("ENV")
	if !found {
		environment = "dev"
	}

	v, err := Configs.NewViper(environment)
	if err != nil {
		return nil, nil, err
	}

	var defaultCfg Configs.DefaultConfig
	if err := v.UnmarshalKey("DefaultConfig", &defaultCfg); err != nil {
		Common.LogError("Failed to unmarshal default config", err, nil)
		// Continue even if default config fails to load completely.
	}

	// Initialize an empty provider config. We'll fill each provider config individually.
	appConfig := &Configs.Config{
		ProviderConfig: &Configs.ProviderConfig{},
		DefaultConfig:  &defaultCfg,
	}

	return v, appConfig, nil
}

// loadProviderConfig unmarshals a provider configuration from Viper for a given key.
func loadProviderConfig(v *viper.Viper, key string) (Configs.ProviderConfig, error) {
	var cfg Configs.ProviderConfig
	err := v.UnmarshalKey(key, &cfg)
	return cfg, err
}

func handleRequest(ctx context.Context) {
	// Load configuration
	v, appCfg, err := LoadConfig()
	if err != nil {
		Common.LogError("Failed to load config", err, nil)
		return
	}

	// Unmarshal provider configurations
	googleCfg, err := loadProviderConfig(v, "googleProvider")
	if err != nil {
		Common.LogError("Failed to unmarshal google trends config", err, nil)
		return
	}

	newsCfg, err := loadProviderConfig(v, "newsProvider")
	if err != nil {
		Common.LogError("Failed to unmarshal news config", err, nil)
		return
	}

	openAICfg, err := loadProviderConfig(v, "openAIProvider")
	if err != nil {
		Common.LogError("Failed to unmarshal openAI config", err, nil)
		return
	}
	// Retrieve secrets
	openAiApiKey, err := GetSecrets(openAICfg.SecretPath, appCfg.DefaultConfig.Region)
	if err != nil {
		Common.LogError("Failed to retrieve openApiKey", err, nil)
		return
	}
	openAICfg.APIKey = openAiApiKey

	// Retrieve secrets (for example, the API key for newsProvider)
	newsApiKey, err := GetSecrets(newsCfg.SecretPath, appCfg.DefaultConfig.Region)
	if err != nil {
		Common.LogError("Failed to retrieve newsApiKey", err, nil)
		return
	}
	newsCfg.APIKey = newsApiKey

	// Initialize shared clients
	metricClient, _ := Metric.NewMetricClient("PollService", appCfg.DefaultConfig.Region)
	s3Client, err := S3.NewS3Client(appCfg.DefaultConfig.S3Bucket, appCfg.DefaultConfig.Region, time.Duration(appCfg.DefaultConfig.DefaultS3Timeout)*time.Second)
	if err != nil {
		Common.LogError("Failed to initialize S3 client", err, nil)
		return
	}

	// Initialize adapters
	trendAdapter := ta.NewTrendsAdapter(Api.NewApiClient(""), metricClient, s3Client)
	//newsAdapter := na.NewNewsAdapter(Api.NewApiClient(newsCfg.APIKey), metricClient)

	// Build complete provider configurations for each provider
	googleAppConfig := &Configs.Config{
		ProviderConfig: &googleCfg,
		DefaultConfig:  appCfg.DefaultConfig,
	}
	//newsAppConfig := &Configs.Config{
	//	ProviderConfig: &newsCfg,
	//	DefaultConfig:  appCfg.DefaultConfig,
	//}
	openAIAAppConfig := &Configs.Config{
		ProviderConfig: &openAICfg,
		DefaultConfig:  appCfg.DefaultConfig,
	}

	// Create providers
	trendProvider := google_trends.NewProvider(googleAppConfig, trendAdapter)
	//newsProvider := news_api.NewProvider(newsAppConfig, newsAdapter)
	// Assume openAI provider uses the same API client (newsApiKey here) if not, adjust accordingly.
	openAiProvider := openAI.NewProvider(openAIAAppConfig, Api.NewApiClient(""), metricClient)

	// Get trends from Google Trends
	trends, err := trendProvider.GetTrends(ctx)
	if err != nil {
		Common.LogError("Failed to get trends", err, nil)
		return
	}

	for _, trend := range trends {
		Common.LogInfo("Trend", map[string]interface{}{"trend": trend})
		//// lets spilt the trend and get the title
		topic := strings.Split(trend, ",")[0]
		Common.LogInfo("Topic", map[string]interface{}{"topic": topic})
		//newsTitles, err := newsProvider.SearchNews(ctx, topic)
		//newsTitles, err := newsProvider.SearchNews(ctx, trend)
		//if err != nil {
		//	Common.LogError("Failed to search news; skipping trend", err, map[string]interface{}{"trend": trend})
		//	continue
		//}
		//Common.LogInfo("News titles", map[string]interface{}{
		//	"trend":      trend,
		//	"newsTitles": newsTitles,
		//})
		//if len(newsTitles) == 0 {
		//	Common.LogInfo("No news titles found", map[string]interface{}{"trend": trend})
		//	continue
		//}
		//
		//// Create a context summary by marshalling news titles
		contextSummaryBytes, err := json.Marshal(trend)
		if err != nil {
			Common.LogError("Failed to marshal news titles", err, map[string]interface{}{"newsTitles": topic})
			continue
		}

		pollQuestions, err := openAiProvider.GeneratePollQuestions(ctx, topic, string(contextSummaryBytes))
		if err != nil {
			Common.LogError("Failed to generate poll questions", err, map[string]interface{}{
				"trend":          trend,
				"contextSummary": string(contextSummaryBytes),
			})
			continue
		}
		Common.LogInfo("Poll questions", map[string]interface{}{
			"trend":         trend,
			"pollQuestions": pollQuestions,
		})
	}
}

func main() {
	lambda.Start(handleRequest)
}
