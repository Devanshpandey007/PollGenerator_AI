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
			want:    []string{"trend1, news Item 1", "trend2, news Item 1, news Item 2, news Item 3"},
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
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/rss?geo=NG", map[string]string{}).Return(
					nil, fmt.Errorf("error"))
			} else {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://api.example.com/trendingsearches/rss?geo=NG", map[string]string{}).Return(
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
