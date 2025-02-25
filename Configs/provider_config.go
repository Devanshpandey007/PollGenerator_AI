package Configs

type ProviderConfig struct {
	ProviderName string
	APIKey       string
	BaseURL      string
	Endpoint     string
	TimeRange    string
	QueryParams  string
	Region       string
	SecretPath   string
	Model        string
}

type DefaultConfig struct {
	// Env configurations
	Environment      string
	Region           string
	DynamoEndpoint   string
	EventTableName   string
	PartyTableName   string
	S3Bucket         string
	DefaultS3Timeout int
	MaxTokens        int
	Temperature      float64
	SearchStartDate  string
}

// Config holds all the configuration settings.
type Config struct {
	ProviderConfig *ProviderConfig
	DefaultConfig  *DefaultConfig
}
