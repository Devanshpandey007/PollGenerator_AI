package Adapters

import (
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Configs"
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
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
		{
			name: "TestNewNewsAdapter_success",
			args: args{
				httpClient:   &Api.ApiClient{},
				metricClient: &Metric.MetricClient{},
			},
			want: &NewsAdapter{
				HttpClient:   &Api.ApiClient{},
				MetricClient: &Metric.MetricClient{},
			},
		},
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
		config         Configs.Config
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
		{
			name: "TestNewsAdapter_SearchNews_success",
			fields: fields{
				HttpClient:   new(Api.MockApiClient),
				MetricClient: new(Metric.MockMetricClient),
			},
			args: args{
				ctx: context.Background(),
				config: Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						BaseURL:  "https://newsapi.org/v2",
						Endpoint: "everything",
						APIKey:   "test_api_key",
					},
					DefaultConfig: &Configs.DefaultConfig{
						SearchStartDate: "2025-01-24",
					},
				},
				topic: "test",
				responseReader: func(res []byte) ([]string, error) {
					return []string{"test"}, nil
				},
			},
			want:    []string{"test"},
			wantErr: false,
		},
		{
			name: "TestNewsAdapter_SearchNews_error",
			fields: fields{
				HttpClient:   new(Api.MockApiClient),
				MetricClient: new(Metric.MockMetricClient),
			},
			args: args{
				ctx: context.Background(),
				config: Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						BaseURL:  "https://newsapi.org/v2",
						Endpoint: "everything",
						APIKey:   "test_api_key",
					},
					DefaultConfig: &Configs.DefaultConfig{
						SearchStartDate: "2025-01-24",
					},
				},
				topic: "test",
				responseReader: func(res []byte) ([]string, error) {
					return nil, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			na := &NewsAdapter{
				HttpClient:   tt.fields.HttpClient,
				MetricClient: tt.fields.MetricClient,
			}
			if tt.wantErr {
				tt.fields.MetricClient.(*Metric.MockMetricClient).On("SendMetric", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://newsapi.org/v2/everything?q=test&from=2025-01-24&sortBy=publishedAt&apiKey=test_api_key", map[string]string{}).Return(nil, fmt.Errorf("error"))
			} else {
				tt.fields.MetricClient.(*Metric.MockMetricClient).On("SendMetric", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://newsapi.org/v2/everything?q=test&from=2025-01-24&sortBy=publishedAt&apiKey=test_api_key", map[string]string{}).Return(&http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("test")),
				}, nil)
			}
			got, err := na.SearchNews(tt.args.ctx, tt.args.config, tt.args.topic, tt.args.responseReader)
			if (err != nil) != tt.wantErr {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://newsapi.org/v2/everything?q=test&from=2025-01-24&sortBy=publishedAt&apiKey=test_api_key", map[string]string{}).Return(nil, fmt.Errorf("error"))
				t.Errorf("SearchNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchNews() got = %v, want %v", got, tt.want)
			}
		})
	}
}
