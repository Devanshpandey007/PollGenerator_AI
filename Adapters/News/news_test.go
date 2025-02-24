package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Configs"
	"context"
	"reflect"
	"testing"
)

func TestNewNewsAdapter(t *testing.T) {
	type args struct {
		httpClient   Api.ApiClientInterface
		metricClient Metric.MetricClientInterface
	}
	tests := []struct {
		name string
		args args
		want *NewsAdapter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNewsAdapter(tt.args.httpClient, tt.args.metricClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNewsAdapter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAdapter_SearchNews(t *testing.T) {
	type fields struct {
		HttpClient   Api.ApiClientInterface
		MetricClient Metric.MetricClientInterface
	}
	type args struct {
		ctx            context.Context
		config         Configs.ProviderConfig
		topic          string
		responseReader func(res []byte) ([]string, error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			na := &NewsAdapter{
				HttpClient:   tt.fields.HttpClient,
				MetricClient: tt.fields.MetricClient,
			}
			got, err := na.SearchNews(tt.args.ctx, tt.args.config, tt.args.topic, tt.args.responseReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchNews() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAdapter_fetchNews(t *testing.T) {
	type fields struct {
		HttpClient   Api.ApiClientInterface
		MetricClient Metric.MetricClientInterface
	}
	type args struct {
		ctx     context.Context
		url     string
		headers map[string]string
		config  Configs.ProviderConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			na := &NewsAdapter{
				HttpClient:   tt.fields.HttpClient,
				MetricClient: tt.fields.MetricClient,
			}
			got, err := na.fetchNews(tt.args.ctx, tt.args.url, tt.args.headers, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchNews() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildNewsUrl(t *testing.T) {
	type args struct {
		config Configs.ProviderConfig
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
			if got := buildNewsUrl(tt.args.config); got != tt.want {
				t.Errorf("buildNewsUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
