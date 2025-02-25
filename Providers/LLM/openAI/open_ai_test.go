package openAI

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Configs"
	"Providers/Providers/LLM"
	"context"
	"reflect"
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
		// TODO: Add test cases.
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
