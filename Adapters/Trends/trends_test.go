package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Configs"
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"reflect"
	"strings"
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
		{
			name: "TestNewTrendsAdapter_success",
			args: args{
				httpClient:   &Api.ApiClient{},
				metricClient: &Metric.MetricClient{},
				s3Client:     &S3.S3Client{},
			},
			want: &TrendsAdapter{
				HttpClient:   &Api.ApiClient{},
				MetricClient: &Metric.MetricClient{},
				S3Client:     &S3.S3Client{},
			},
		},
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
		{
			name: "TestTrendsAdapter_GetTrends_success",
			fields: fields{
				HttpClient:   new(Api.MockApiClient),
				S3Client:     new(S3.MockS3Client),
				MetricClient: new(Metric.MockMetricClient),
			},
			args: args{
				ctx: context.Background(),
				providerConfigs: []Configs.ProviderConfig{
					{
						BaseURL:   "https://api.example.com",
						Endpoint:  "trendingsearches",
						TimeRange: "daily",
						Region:    "NG",
					},
				},
				responseReaderFunc: func(res []byte) ([]string, error) {
					// return a list of trends
					return []string{"trend1", "trend2"}, nil
				},
			},
			want:    []string{"trend1", "trend2"},
			wantErr: false,
		},
		{
			name: "TestTrendsAdapter_GetTrends_error",
			fields: fields{
				HttpClient:   new(Api.MockApiClient),
				S3Client:     new(S3.MockS3Client),
				MetricClient: new(Metric.MockMetricClient),
			},
			args: args{
				ctx: context.Background(),
				providerConfigs: []Configs.ProviderConfig{
					{
						BaseURL:   "https://api.example.com",
						Endpoint:  "trendingsearches",
						TimeRange: "daily",
						Region:    "NG",
					},
				},
				responseReaderFunc: func(res []byte) ([]string, error) {
					return nil, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &TrendsAdapter{
				HttpClient:   tt.fields.HttpClient,
				S3Client:     tt.fields.S3Client,
				MetricClient: tt.fields.MetricClient,
			}
			tt.fields.MetricClient.(*Metric.MockMetricClient).On("SendMetric", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			if tt.wantErr {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/daily?geo=NG", map[string]string{}).Return(
					nil, fmt.Errorf("error"))
			} else {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/daily?geo=NG", map[string]string{}).Return(
					&http.Response{
						Body:       io.NopCloser(strings.NewReader(`{"trends": ["trend1", "trend2"]}`)),
						StatusCode: 200,
					}, nil)
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
