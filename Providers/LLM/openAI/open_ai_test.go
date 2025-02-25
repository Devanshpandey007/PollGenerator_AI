package openAI

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Configs"
	"Providers/Providers/LLM"
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		config       *Configs.Config
		httpClient   Api.ApiClientInterface
		metricClient Metric.MetricClientInterface
	}
	tests := []struct {
		name string
		args args
		want *OpenAIProvider
	}{
		{
			name: "TestNewProvider_success",
			args: args{
				config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "OpenAIProvider",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
				httpClient:   new(Api.MockApiClient),
				metricClient: new(Metric.MockMetricClient),
			},
			want: &OpenAIProvider{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "OpenAIProvider",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
				HttpClient:   new(Api.MockApiClient),
				MetricClient: new(Metric.MockMetricClient),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProvider(tt.args.config, tt.args.httpClient, tt.args.metricClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenAIProvider_GeneratePollQuestions(t *testing.T) {
	type fields struct {
		IProviderLLM LLM.IProviderLLM
		Config       *Configs.Config
		HttpClient   Api.ApiClientInterface
		MetricClient Metric.MetricClientInterface
	}
	type args struct {
		ctx            context.Context
		topic          string
		contextSummary string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OpenAIProvider{
				IProviderLLM: tt.fields.IProviderLLM,
				Config:       tt.fields.Config,
				HttpClient:   tt.fields.HttpClient,
				MetricClient: tt.fields.MetricClient,
			}
			got, err := p.GeneratePollQuestions(tt.args.ctx, tt.args.topic, tt.args.contextSummary)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePollQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GeneratePollQuestions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenAIProvider_GetConfig(t *testing.T) {
	type fields struct {
		IProviderLLM LLM.IProviderLLM
		Config       *Configs.Config
		HttpClient   Api.ApiClientInterface
		MetricClient Metric.MetricClientInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   Configs.ProviderConfig
	}{
		{
			name: "TestGetConfig_success",
			fields: fields{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "OpenAIProvider",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
			},
			want: Configs.ProviderConfig{
				ProviderName: "OpenAIProvider",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OpenAIProvider{
				IProviderLLM: tt.fields.IProviderLLM,
				Config:       tt.fields.Config,
				HttpClient:   tt.fields.HttpClient,
				MetricClient: tt.fields.MetricClient,
			}
			if got := p.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenAIProvider_GetProviderName(t *testing.T) {
	type fields struct {
		IProviderLLM LLM.IProviderLLM
		Config       *Configs.Config
		HttpClient   Api.ApiClientInterface
		MetricClient Metric.MetricClientInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OpenAIProvider{
				IProviderLLM: tt.fields.IProviderLLM,
				Config:       tt.fields.Config,
				HttpClient:   tt.fields.HttpClient,
				MetricClient: tt.fields.MetricClient,
			}
			if got := p.GetProviderName(); got != tt.want {
				t.Errorf("GetProviderName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenAIProvider_makePromptRequest(t *testing.T) {
	type fields struct {
		IProviderLLM LLM.IProviderLLM
		Config       *Configs.Config
		HttpClient   Api.ApiClientInterface
		MetricClient Metric.MetricClientInterface
	}
	type args struct {
		ctx    context.Context
		prompt string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestOpenAIProvider_makePromptRequest_success",
			fields: fields{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "OpenAIProvider",
						BaseURL:      "https://api.openai.com/v1",
						Endpoint:     "chat/completions",
						Model:        "gpt-3.5-turbo",
						APIKey:       "key",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
						MaxTokens:   100,
						Temperature: 0.7,
					},
				},
				HttpClient:   new(Api.MockApiClient),
				MetricClient: new(Metric.MockMetricClient),
			},
			args: args{
				ctx:    context.Background(),
				prompt: "prompt",
			},
			want:    "response",
			wantErr: false,
		},
		{
			name: "TestOpenAIProvider_makePromptRequest_error",
			fields: fields{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "OpenAIProvider",
						BaseURL:      "https://api.openai.com/v1",
						Endpoint:     "chat/completions",
						Model:        "gpt-3.5-turbo",
						APIKey:       "key",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
						MaxTokens:   100,
						Temperature: 0.7,
					},
				},
				HttpClient:   new(Api.MockApiClient),
				MetricClient: new(Metric.MockMetricClient),
			},
			args: args{
				ctx:    context.Background(),
				prompt: "prompt",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OpenAIProvider{
				IProviderLLM: tt.fields.IProviderLLM,
				Config:       tt.fields.Config,
				HttpClient:   tt.fields.HttpClient,
				MetricClient: tt.fields.MetricClient,
			}
			tt.fields.MetricClient.(*Metric.MockMetricClient).On("SendMetric", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			if tt.wantErr {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakePostRequest", tt.args.ctx, "https://api.openai.com/v1/chat/completions", map[string]string{
					"Authorization": "Bearer key",
					"Content-Type":  "application/json",
				}, mock.Anything).Return(
					nil, fmt.Errorf("error"))
			} else {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakePostRequest", tt.args.ctx, "https://api.openai.com/v1/chat/completions", map[string]string{
					"Authorization": "Bearer key",
					"Content-Type":  "application/json",
				}, mock.Anything).Return(
					&http.Response{
						Body: io.NopCloser(strings.NewReader(`{
    "id": "chatcmpl-B4ybHIzaCb6nMqZUeEpkvQWM5oGOg",
    "object": "chat.completion",
    "created": 1740525959,
    "model": "gpt-4o-mini-2024-07-18",
    "choices": [
        {
            "index": 0,
            "message": {
                "role": "assistant",
                "content": "response",
                "refusal": null
            },
            "logprobs": null,
            "finish_reason": "stop"
        }
    ],
    "usage": {
        "prompt_tokens": 130,
        "completion_tokens": 22,
        "total_tokens": 152,
        "prompt_tokens_details": {
            "cached_tokens": 0,
            "audio_tokens": 0
        },
        "completion_tokens_details": {
            "reasoning_tokens": 0,
            "audio_tokens": 0,
            "accepted_prediction_tokens": 0,
            "rejected_prediction_tokens": 0
        }
    },
    "service_tier": "default",
    "system_fingerprint": "fp_709714d124"
}`)),
						StatusCode: 200,
					}, nil)
			}
			got, err := p.makePromptRequest(tt.args.ctx, tt.args.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("makePromptRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makePromptRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildPromptUrl(t *testing.T) {
	type args struct {
		config *Configs.ProviderConfig
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildPromptUrl(tt.args.config); got != tt.want {
				t.Errorf("buildPromptUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
