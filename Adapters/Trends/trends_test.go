package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Configs"
	"context"
	"reflect"
	"testing"
)

func TestNewTrendsAdapter(t *testing.T) {
	type args struct {
		httpClient   Api.ApiClientInterface
		metricClient Metric.MetricClientInterface
		s3Client     S3.S3ClientInterface
	}
	tests := []struct {
		name string
		args args
		want *TrendsAdapter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrendsAdapter(tt.args.httpClient, tt.args.metricClient, tt.args.s3Client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrendsAdapter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrendsAdapter_GetTrends(t *testing.T) {
	type fields struct {
		HttpClient   Api.ApiClientInterface
		S3Client     S3.S3ClientInterface
		MetricClient Metric.MetricClientInterface
	}
	type args struct {
		ctx                context.Context
		providerConfigs    []Configs.ProviderConfig
		responseReaderFunc func([]byte) ([]string, error)
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
			ta := &TrendsAdapter{
				HttpClient:   tt.fields.HttpClient,
				S3Client:     tt.fields.S3Client,
				MetricClient: tt.fields.MetricClient,
			}
			got, err := ta.GetTrends(tt.args.ctx, tt.args.providerConfigs, tt.args.responseReaderFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrends() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTrends() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrendsAdapter_fetchTrends(t *testing.T) {
	type fields struct {
		HttpClient   Api.ApiClientInterface
		S3Client     S3.S3ClientInterface
		MetricClient Metric.MetricClientInterface
	}
	type args struct {
		ctx    context.Context
		config Configs.ProviderConfig
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
			ta := &TrendsAdapter{
				HttpClient:   tt.fields.HttpClient,
				S3Client:     tt.fields.S3Client,
				MetricClient: tt.fields.MetricClient,
			}
			got, err := ta.fetchTrends(tt.args.ctx, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchTrends() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchTrends() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildTrendsUrl(t *testing.T) {
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
			if got := buildTrendsUrl(tt.args.config); got != tt.want {
				t.Errorf("buildTrendsUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
