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
			want:    []string{"trend1, news Item 1", "trend2, news Item 1, news Item 2, news Item 3"},
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
				tt.fields.Adapter.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/rss?geo=NG", map[string]string{}).Return(
					nil, fmt.Errorf("error"))
			} else {
				tt.fields.Adapter.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/rss?geo=NG", map[string]string{}).Return(
					&http.Response{
						Body: io.NopCloser(strings.NewReader(`
							<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:ht="https://trends.google.com/trending/rss" version="2.0">
								<channel>
									<title>Daily Search Trends</title>
									<description>Recent searches</description>
									<link>https://trends.google.com/trending/rss?geo=US</link>
									<atom:link href="https://trends.google.com/trending/rss?geo=US" rel="self" type="application/rss+xml"/>
									<item>
										<title>trend1</title>
										<ht:approx_traffic>2000+</ht:approx_traffic>
										<description/>
										<link>https://trends.google.com/trending/rss?geo=US</link>
										<pubDate>Tue, 25 Feb 2025 13:20:00 -0800</pubDate>
										<ht:picture>https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcSsGdm1ftf6kSRV_6TDut4L9IUtqC4-Agf_K5zwrPWlCH0wbc6hHFRZDDlsuRY</ht:picture>
										<ht:picture_source>Yahoo Finance</ht:picture_source>
										<ht:news_item>
											<ht:news_item_title>news Item 1</ht:news_item_title>
											<ht:news_item_snippet/>
											<ht:news_item_url>https://finance.yahoo.com/news/super-micro-stock-drops-as-nasdaq-deadline-to-avoid-delisting-approaches-200716294.html</ht:news_item_url>
											<ht:news_item_picture>https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcSsGdm1ftf6kSRV_6TDut4L9IUtqC4-Agf_K5zwrPWlCH0wbc6hHFRZDDlsuRY</ht:news_item_picture>
											<ht:news_item_source>Yahoo Finance</ht:news_item_source>
										</ht:news_item>
									</item>
									<item>
										<title>trend2</title>
										<ht:approx_traffic>2000+</ht:approx_traffic>
										<description/>
										<link>https://trends.google.com/trending/rss?geo=US</link>
										<pubDate>Tue, 25 Feb 2025 13:10:00 -0800</pubDate>
										<ht:picture>https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcRpAv-AFCVt8yJkF9l4fO1B0MCTVsxk0xvd5FUaKi4Akmy3RGy-2s_AGSx11zI</ht:picture>
										<ht:picture_source>CNN</ht:picture_source>
										<ht:news_item>
											<ht:news_item_title>news Item 1</ht:news_item_title>
											<ht:news_item_snippet/>
											<ht:news_item_url>https://www.cnn.com/2025/02/25/health/congo-mystery-illness/index.html</ht:news_item_url>
											<ht:news_item_picture>https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcRpAv-AFCVt8yJkF9l4fO1B0MCTVsxk0xvd5FUaKi4Akmy3RGy-2s_AGSx11zI</ht:news_item_picture>
											<ht:news_item_source>CNN</ht:news_item_source>
										</ht:news_item>
										<ht:news_item>
											<ht:news_item_title>news Item 2</ht:news_item_title>
											<ht:news_item_snippet/>
											<ht:news_item_url>https://www.livescience.com/health/viruses-infections-disease/unidentified-illnesses-have-killed-over-50-people-in-congo-in-last-5-weeks-who-reports</ht:news_item_url>
											<ht:news_item_picture>https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcQtuhaGTft9STrOxvBRauOaga_ug_niXA-ysnyX0MqiLDf7xGmLhKNj-khRi8I</ht:news_item_picture>
											<ht:news_item_source>Live Science</ht:news_item_source>
										</ht:news_item>
										<ht:news_item>
											<ht:news_item_title>news Item 3</ht:news_item_title>
											<ht:news_item_snippet/>
											<ht:news_item_url>https://www.washingtonpost.com/world/2025/02/25/unknown-illness-hemorrhagic-fever-congo-africa/</ht:news_item_url>
											<ht:news_item_picture>https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcTIdGG3DFHF2GnN7_RzYy1tP8Sp7YNd-Zeg0UW9DzR1zwtZO3JiGkoEPXc4B9w</ht:news_item_picture>
											<ht:news_item_source>The Washington Post</ht:news_item_source>
										</ht:news_item>
									</item>
								</channel>
							</rss>`)),
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
