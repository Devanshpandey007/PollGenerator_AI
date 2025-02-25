package google_trends

import (
	Adapters "Providers/Adapters/Trends"
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Clients/S3"
	"Providers/Configs"
	"Providers/Providers/Trends"
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestGoogleTrendsProvider_GetConfig(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Config          *Configs.Config
		Adapter         *Adapters.TrendsAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   Configs.ProviderConfig
	}{
		{
			name: "TestGoogleTrendsProvider_GetConfig_success",
			fields: fields{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "GoogleTrendsProvider",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
			},
			want: Configs.ProviderConfig{
				ProviderName: "GoogleTrendsProvider",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Config,
				Adapter:         tt.fields.Adapter,
			}
			if got := p.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleTrendsProvider_GetTrends(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Config          *Configs.Config
		Adapter         *Adapters.TrendsAdapter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "TestGoogleTrendsProvider_GetTrends_success",
			fields: fields{
				Adapter: &Adapters.TrendsAdapter{
					HttpClient:   new(Api.MockApiClient),
					S3Client:     new(S3.MockS3Client),
					MetricClient: new(Metric.MockMetricClient),
				},
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "GoogleTrendsProvider",
						BaseURL:      "https://api.example.com",
						Endpoint:     "trendingsearches",
						TimeRange:    "daily",
						Region:       "NG",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []string{"Super Bowl 2024", "Bitcoin Price Surge", "Ukraine Conflict"},
			wantErr: false,
		},
		{
			name: "TestGoogleTrendsProvider_GetTrends_error",
			fields: fields{
				Adapter: &Adapters.TrendsAdapter{
					HttpClient:   new(Api.MockApiClient),
					S3Client:     new(S3.MockS3Client),
					MetricClient: new(Metric.MockMetricClient),
				},
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "GoogleTrendsProvider",
						BaseURL:      "https://api.example.com",
						Endpoint:     "trendingsearches",
						TimeRange:    "daily",
						Region:       "NG",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Config,
				Adapter:         tt.fields.Adapter,
			}
			tt.fields.Adapter.MetricClient.(*Metric.MockMetricClient).On("SendMetric", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			if tt.wantErr {
				tt.fields.Adapter.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/daily?geo=NG", map[string]string{}).Return(
					nil, fmt.Errorf("error"))
			} else {
				tt.fields.Adapter.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/daily?geo=NG", map[string]string{}).Return(
					&http.Response{
						Body: io.NopCloser(strings.NewReader(`<html>
								  <head>
									<title>Google Trends - Daily Searches</title>
								  </head>
								  <body>
									<div class="summary-text">Super Bowl 2024</div>
									<div class="summary-text">Bitcoin Price Surge</div>
									<div class="summary-text">Ukraine Conflict</div>
								  </body>
								</html>
								`)),
						StatusCode: 200,
					}, nil)
			}
			got, err := p.GetTrends(tt.args.ctx)
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

func TestGoogleTrendsProvider_ResponseReaderFunc(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Config          *Configs.Config
		Adapter         *Adapters.TrendsAdapter
	}
	type args struct {
		body []byte
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
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Config,
				Adapter:         tt.fields.Adapter,
			}
			got, err := p.ResponseReaderFunc(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResponseReaderFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResponseReaderFunc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleTrendsProvider_SaveTrend(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Config          *Configs.Config
		Adapter         *Adapters.TrendsAdapter
	}
	type args struct {
		ctx   context.Context
		trend string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Config,
				Adapter:         tt.fields.Adapter,
			}
			if err := p.SaveTrend(tt.args.ctx, tt.args.trend); (err != nil) != tt.wantErr {
				t.Errorf("SaveTrend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoogleTrendsProvider_SaveTrendForRegion(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Config          *Configs.Config
		Adapter         *Adapters.TrendsAdapter
	}
	type args struct {
		ctx    context.Context
		trend  string
		region string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Config,
				Adapter:         tt.fields.Adapter,
			}
			if err := p.SaveTrendForRegion(tt.args.ctx, tt.args.trend, tt.args.region); (err != nil) != tt.wantErr {
				t.Errorf("SaveTrendForRegion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewProvider(t *testing.T) {
	type args struct {
		config *Configs.Config
		trends *Adapters.TrendsAdapter
	}
	tests := []struct {
		name string
		args args
		want *GoogleTrendsProvider
	}{
		{
			name: "TestNewProvider_success",
			args: args{
				config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "GoogleTrendsProvider",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
			},
			want: &GoogleTrendsProvider{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "GoogleTrendsProvider",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment: "test",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProvider(tt.args.config, tt.args.trends); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
