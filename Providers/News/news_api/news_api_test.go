package news_api

import (
	Adapters "Providers/Adapters/News"
	Api "Providers/Clients/Http"
	"Providers/Clients/Metric"
	"Providers/Configs"
	"Providers/Providers/News"
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		config *Configs.Config
		news   *Adapters.NewsAdapter
	}
	tests := []struct {
		name string
		args args
		want *NewsApiProvider
	}{
		{
			name: "TestNewProvider_success",
			args: args{
				config: &Configs.Config{},
				news:   &Adapters.NewsAdapter{},
			},
			want: &NewsApiProvider{
				Config:  &Configs.Config{},
				Adapter: &Adapters.NewsAdapter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProvider(tt.args.config, tt.args.news); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsApiProvider_ResponseReaderFunc(t *testing.T) {
	type fields struct {
		IProviderNews News.IProviderNews
		Config        *Configs.Config
		Adapter       *Adapters.NewsAdapter
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
			p := &NewsApiProvider{
				IProviderNews: tt.fields.IProviderNews,
				Config:        tt.fields.Config,
				Adapter:       tt.fields.Adapter,
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

func TestNewsApiProvider_SearchNews(t *testing.T) {
	type fields struct {
		IProviderNews News.IProviderNews
		Config        *Configs.Config
		Adapter       *Adapters.NewsAdapter
	}
	type args struct {
		ctx   context.Context
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []News.NewsArticle
		wantErr bool
	}{
		{
			name: "TestNewsApiProvider_SearchNews_success",
			fields: fields{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "newsapi",
						BaseURL:      "https://newsapi.org/v2",
						Endpoint:     "everything",
						APIKey:       "test_api_key",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment:     "test",
						SearchStartDate: "2025-01-24",
					},
				},
				Adapter: &Adapters.NewsAdapter{
					HttpClient:   new(Api.MockApiClient),
					MetricClient: new(Metric.MockMetricClient),
				},
			},
			args: args{
				ctx:   context.Background(),
				query: "test",
			},
			want: []News.NewsArticle{
				{
					Title:       "Nvidia Earnings, Inflation: What to Watch This Week",
					Description: "Last, but not least. The final member of the so-called Magnificent Seven tech stocks to report in this earnings season is Nvidia, which will update investors late Wednesday.Those results, from the semiconductor company at the center of the artificial-intellig…",
					PublishedAt: "2025-02-24T08:13:11Z",
					Content:     "Last, but not least. The final member of the so-called Magnificent Seven tech stocks to report in this earnings season is Nvidia, which will update investors late Wednesday.\r\nThose results, from the … [+265 chars]",
				},
			},
			wantErr: false,
		},
		{
			name: "TestNewsApiProvider_SearchNews_error",
			fields: fields{
				Config: &Configs.Config{
					ProviderConfig: &Configs.ProviderConfig{
						ProviderName: "newsapi",
						BaseURL:      "https://newsapi.org/v2",
						Endpoint:     "everything",
						APIKey:       "test_api_key",
					},
					DefaultConfig: &Configs.DefaultConfig{
						Environment:     "test",
						SearchStartDate: "2025-01-24",
					},
				},
				Adapter: &Adapters.NewsAdapter{
					HttpClient:   new(Api.MockApiClient),
					MetricClient: new(Metric.MockMetricClient),
				},
			},
			args: args{
				ctx:   context.Background(),
				query: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &NewsApiProvider{
				IProviderNews: tt.fields.IProviderNews,
				Config:        tt.fields.Config,
				Adapter:       tt.fields.Adapter,
			}
			tt.fields.Adapter.MetricClient.(*Metric.MockMetricClient).On("SendMetric", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			if tt.wantErr {
				tt.fields.Adapter.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://newsapi.org/v2/everything?q=test&from=2025-01-24&sortBy=publishedAt&apiKey=test_api_key", map[string]string{}).Return(nil, fmt.Errorf("error"))
			} else {
				tt.fields.Adapter.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://newsapi.org/v2/everything?q=test&from=2025-01-24&sortBy=publishedAt&apiKey=test_api_key", map[string]string{}).Return(&http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(strings.NewReader(`{
								"status": "ok",
								"totalResults": 5,
								"articles": [
								{
								"source": {
								"id": "the-wall-street-journal",
								"name": "The Wall Street Journal"
								},
								"author": "WSJ Staff",
								"title": "Nvidia Earnings, Inflation: What to Watch This Week",
								"description": "Last, but not least. The final member of the so-called Magnificent Seven tech stocks to report in this earnings season is Nvidia, which will update investors late Wednesday.Those results, from the semiconductor company at the center of the artificial-intellig…",
								"url": "https://www.wsj.com/livecoverage/stock-market-today-dow-sp500-nasdaq-02-24-2025/card/nvidia-earnings-inflation-what-to-watch-this-week-WOLXuHPDYf6Pmtjvzmw7",
								"urlToImage": "https://images.wsj.net/im-925349/social",
								"publishedAt": "2025-02-24T08:13:11Z",
								"content": "Last, but not least. The final member of the so-called Magnificent Seven tech stocks to report in this earnings season is Nvidia, which will update investors late Wednesday.\r\nThose results, from the … [+265 chars]"
								}
								]
								}`)),
				}, nil)
			}
			got, err := p.SearchNews(tt.args.ctx, tt.args.query)
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
