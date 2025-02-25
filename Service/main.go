package main

import (
	na "Providers/Adapters/News"
	ta "Providers/Adapters/Trends"
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Common"
	"Providers/Configs"
	"Providers/Providers/LLM/openAI"
	"Providers/Providers/News/news_api"
	"Providers/Providers/Trends/google_trends"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"time"
)

// LoadConfig loads configuration from environment variables via Config
func LoadConfig() (*Configs.Config, error) {
	environment, foundEnv := os.LookupEnv("ENV")
	if !foundEnv {
		environment = "test"
	}

	v, err := Configs.NewViper(environment)
	if err != nil {
		return nil, err
	}

	// Unmarshal the configuration from Config
	config := &Configs.Config{}
	err = v.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func handleRequest(ctx context.Context) {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	// Get the api key for the news API
	secrets, err := GetSecrets(config.ProviderConfig.SecretPath, config.DefaultConfig.Region)
	if err != nil {
		Common.LogError("Failed to retrieve secrets", err, nil)
		return
	}
	config.ProviderConfig.APIKey = secrets

	// Initialize the clients
	metricClient, _ := Metric.NewMetricClient("PollService", config.DefaultConfig.Region)
	s3Client, err := S3.NewS3Client(config.DefaultConfig.S3Bucket, config.DefaultConfig.Region, time.Duration(config.DefaultConfig.DefaultS3Timeout)*time.Second)

	// Initialize the adapter
	trendAdapter := ta.NewTrendsAdapter(Api.NewApiClient(""), metricClient, s3Client)
	newsAdapter := na.NewNewsAdapter(Api.NewApiClient(secrets), metricClient)

	// Create the providers
	trendProvider := google_trends.NewProvider(config, trendAdapter)
	newsProvider := news_api.NewProvider(config, newsAdapter)
	openAiProvider := openAI.NewProvider(config, Api.NewApiClient(secrets), metricClient)

	trends, err := trendProvider.GetTrends(ctx)
	if err != nil {
		return
	}
	for _, trend := range trends {
		Common.LogInfo("Trend", map[string]interface{}{
			"trend": trend,
		})
		newsTitles, err := newsProvider.SearchNews(ctx, trend)
		if err != nil {
			Common.LogError("Failed to search news ... Skipping provider", err, map[string]interface{}{
				"trend": trend,
			})
			continue
		}
		Common.LogInfo("News titles", map[string]interface{}{
			"trend":      trend,
			"newsTitles": newsTitles,
		})

		// Create a context summary by joining article titles.
		contextSummary, err := json.Marshal(newsTitles)
		if err != nil {
			Common.LogError("Failed to marshal news titles", err, map[string]interface{}{
				"newsTitles": newsTitles,
			})
			continue
		}
		pollQuestions, err := openAiProvider.GeneratePollQuestions(ctx, trend, string(contextSummary))
		if err != nil {
			Common.LogError("Failed to generate poll questions", err, map[string]interface{}{
				"trend":          trend,
				"contextSummary": contextSummary,
			})
			return
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
