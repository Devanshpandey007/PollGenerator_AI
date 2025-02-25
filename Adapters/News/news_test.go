package Adapters

import (
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
		responseReader func(res []byte) ([]News.NewsArticle, error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []News.NewsArticle
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
				responseReader: func(res []byte) ([]News.NewsArticle, error) {
					return []News.NewsArticle{
						{
							Title:       "test",
							Description: "test",
							PublishedAt: "2025-01-24",
							Content:     "foo bar",
						},
					}, nil
				},
			},
			want: []News.NewsArticle{
				{
					Title:       "test",
					Description: "test",
					PublishedAt: "2025-01-24",
					Content:     "foo bar",
				},
			},
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
				responseReader: func(res []byte) ([]News.NewsArticle, error) {
					return []News.NewsArticle{}, nil
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
			tt.fields.MetricClient.(*Metric.MockMetricClient).On("SendMetric", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			if tt.wantErr {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://newsapi.org/v2/everything?q=test&from=2025-01-24&sortBy=publishedAt&apiKey=test_api_key", map[string]string{}).Return(nil, fmt.Errorf("error"))
			} else {
				tt.fields.HttpClient.(*Api.MockApiClient).On("MakeGetRequest", tt.args.ctx, "https://newsapi.org/v2/everything?q=test&from=2025-01-24&sortBy=publishedAt&apiKey=test_api_key", map[string]string{}).Return(&http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(strings.NewReader(`{
								"status": "ok",
								"totalResults": 5,
								-"articles": [
								-{
								-"source": {
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
								},
								-{
								-"source": {
								"id": "the-wall-street-journal",
								"name": "The Wall Street Journal"
								},
								"author": "John Gruber",
								"title": "Barbara Broccoli on Amazon, in Private: ‘These People Are Fucking Idiots’",
								"description": null,
								"url": "https://www.wsj.com/business/media/james-bond-movies-amazon-barbara-broccoli-0b04f0db?st=P86jBM&reflink=desktopwebshare_permalink",
								"urlToImage": null,
								"publishedAt": "2025-02-20T19:42:32Z",
								"content": "Erich Schwartzel and Jessica Toonkel, reporting for The Wall Street Journal back on December 19, under the headline “Where Is James Bond? Trapped in an Ugly Stalemate With Amazon” (News+ link):\n\n\n Ne… [+2626 chars]"
								},
								-{
								-"source": {
								"id": "the-wall-street-journal",
								"name": "The Wall Street Journal"
								},
								"author": "Jeanne Whalen and Justin Lahart",
								"title": "Hiring Slows but Remains Solid, With Economy Adding 143,000 Jobs",
								"description": "The gain in jobs was lower than expected, but the job counts for November and December were revised upward by a combined 100,000.",
								"url": "https://www.wsj.com/economy/jobs/jobs-report-january-2025-unemployment-economy-1cb95d5b",
								"urlToImage": "https://s.yimg.com/ny/api/res/1.2/CH1j_ZHogwnCx7wxXwHLMw--/YXBwaWQ9aGlnaGxhbmRlcjt3PTEyMDA7aD02MDA-/https://media.zenfs.com/en/the_wall_street_journal_hosted_996/5c73169e84b7400bfd31cb49c906ab00",
								"publishedAt": "2025-02-07T21:22:00Z",
								"content": "The job market kept chugging along in January, albeit at a slower pace than the previous two months.\r\nThe U.S. economy added 143,000 jobs last month and the unemployment rate edged down to 4%, the La… [+6683 chars]"
								},
								-{
								-"source": {
								"id": "the-wall-street-journal",
								"name": "The Wall Street Journal"
								},
								"author": "Heard Editors",
								"title": "Heard on the Street Thursday Recap: Chip Spending",
								"description": "It is never a good thing when your biggest customer breaks away.That seems to be the case for UPS, which said it was cutting volumes from its largest customer—Amazon—by half. UPS warned that revenue would come in about $6 billion lower than expected this year…",
								"url": "https://www.wsj.com/livecoverage/stock-market-today-pce-inflation-data-01-31-2025/card/heard-on-the-street-thursday-recap-chip-spending-hzCLFa49tgfXYAIFywYw",
								"urlToImage": "https://images.wsj.net/im-56647844",
								"publishedAt": "2025-01-31T08:02:42Z",
								"content": "It is never a good thing when your biggest customer breaks away.\r\nThat seems to be the case for UPS, which said it was cutting volumes from its largest customerAmazonby half. UPS warned that revenue … [+457 chars]"
								},
								-{
								-"source": {
								"id": "the-wall-street-journal",
								"name": "The Wall Street Journal"
								},
								"author": "www.wsj.com",
								"title": "She Made Orgasmic Meditation Her Life. Not Even Prison Will Stop Her...",
								"description": "Nicole Daedone’s supporters call her a ‘visionary.’ Now facing a trial for allegedly exploiting employees, she is out to preserve her legacy. Will orgasmic meditation have a future?",
								"url": "https://www.wsj.com/style/nicole-daedone-one-taste-founder-prison-e1dc5eff",
								"urlToImage": "https://images.wsj.net/im-18303885/social",
								"publishedAt": "2025-01-25T16:00:03Z",
								"content": "For over a decade, Nicole Daedone presided over an unusual wellness empire that transformed sexual stimulation into a meditation practice meant to empower women. Through the 2010s, as the founder of … [+591 chars]"
								}
								]
								}`)),
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
