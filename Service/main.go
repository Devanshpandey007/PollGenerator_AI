package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/spf13/viper"

	ta "Providers/Adapters/Trends"
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Common"
	"Providers/Configs"
	openAIProvider "Providers/Providers/LLM/openAI"
	googleTrends "Providers/Providers/Trends/google_trends"
)

// LoadConfig initializes Viper, loads your default config, and prepares
// a Configs.Config struct. Simplified error handling below.
func LoadConfig() (*viper.Viper, *Configs.Config, error) {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "dev"
	}

	v, err := Configs.NewViper(environment)
	if err != nil {
		return nil, nil, err
	}

	var defaultCfg Configs.DefaultConfig
	if err := v.UnmarshalKey("DefaultConfig", &defaultCfg); err != nil {
		// Log, but continue with partially loaded config
		Common.LogError("Failed to unmarshal default config", err, nil)
	}

	appConfig := &Configs.Config{
		ProviderConfig: &Configs.ProviderConfig{},
		DefaultConfig:  &defaultCfg,
	}

	return v, appConfig, nil
}

// loadProviderConfig extracts a ProviderConfig by key. If missing or invalid,
// returns an error.
func loadProviderConfig(v *viper.Viper, key string) (Configs.ProviderConfig, error) {
	var cfg Configs.ProviderConfig
	err := v.UnmarshalKey(key, &cfg)
	return cfg, err
}

// handleRequest is your main AWS Lambda entry point.
func handleRequest(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Load config
	v, appCfg, err := LoadConfig()
	if err != nil {
		Common.LogError("Failed to load config", err, nil)
		return serverError(err), err
	}

	// 2. Unmarshal each provider’s config
	googleCfg, err := loadProviderConfig(v, "googleProvider")
	if err != nil {
		return logAndReturnError("Failed to unmarshal google trends config", err)
	}

	openAICfg, err := loadProviderConfig(v, "openAIProvider")
	if err != nil {
		return logAndReturnError("Failed to unmarshal openAI config", err)
	}

	// 3. Retrieve secrets (OpenAI key, etc.)
	openAiApiKey, err := GetSecrets(openAICfg.SecretPath, appCfg.DefaultConfig.Region)
	if err != nil {
		return logAndReturnError("Failed to retrieve openApiKey", err)
	}
	openAICfg.APIKey = openAiApiKey

	// 4. Initialize shared clients
	metricClient, _ := Metric.NewMetricClient("PollService", appCfg.DefaultConfig.Region)
	s3Client, err := S3.NewS3Client(
		appCfg.DefaultConfig.S3Bucket,
		appCfg.DefaultConfig.Region,
		time.Duration(appCfg.DefaultConfig.DefaultS3Timeout)*time.Second,
	)
	if err != nil {
		return logAndReturnError("Failed to initialize S3 client", err)
	}

	// 5. Initialize your adapters & providers
	trendAdapter := ta.NewTrendsAdapter(Api.NewApiClient(""), metricClient, s3Client)
	//newsAdapter := na.NewNewsAdapter(Api.NewApiClient(newsCfg.APIKey), metricClient)

	googleAppConfig := &Configs.Config{
		ProviderConfig: &googleCfg,
		DefaultConfig:  appCfg.DefaultConfig,
	}
	openAIAppConfig := &Configs.Config{
		ProviderConfig: &openAICfg,
		DefaultConfig:  appCfg.DefaultConfig,
	}

	trendProvider := googleTrends.NewProvider(googleAppConfig, trendAdapter)
	openAiProvider := openAIProvider.NewProvider(openAIAppConfig, Api.NewApiClient(""), metricClient)

	// 6. Get trends from Google Trends
	trends, err := trendProvider.GetTrends(ctx)
	if err != nil {
		return logAndReturnError("Failed to get trends", err)
	}

	var questions []string

	for _, trend := range trends {
		// If the trend has multiple comma parts, we only take the first as "topic".
		topic := strings.Split(trend, ",")[0]
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
		Common.LogInfo("Processing Trend", map[string]interface{}{
			"original": trend,
			"topic":    topic,
		})

		// Create a “context summary” from the trend itself (in future, from news titles)
		contextSummary, err := json.Marshal(trend)
		if err != nil {
			Common.LogError("Failed to marshal trend context", err, nil)
			continue
		}

		pollQuestions, err := openAiProvider.GeneratePollQuestions(ctx, topic, string(contextSummary))
		if err != nil {
			Common.LogError("Failed to generate poll questions", err, map[string]interface{}{
				"trend":          trend,
				"contextSummary": string(contextSummary),
			})
			continue
		}

		Common.LogInfo("Poll questions generated", map[string]interface{}{
			"trend":         trend,
			"pollQuestions": pollQuestions,
		})
		questions = append(questions, pollQuestions)
	}

	// 7. Combine final results & return
	questionsString := strings.Join(questions, ", ")
	return successResponse(questionsString), nil
}

func main() {
	lambda.Start(handleRequest)
}

/* ---------------------------------------------------
   Helper Functions
--------------------------------------------------- */

// successResponse returns a 200 with a JSON body
func successResponse(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf(`{"message":"%s"}`, msg),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

// serverError returns a 500 with a JSON body
func serverError(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

// logAndReturnError logs the error and returns a standardized 500 response
func logAndReturnError(logMessage string, err error) (events.APIGatewayProxyResponse, error) {
	Common.LogError(logMessage, err, nil)
	return serverError(err), err
}
