package news_api

import (
	Adapters "Providers/Adapters/News"
	"Providers/Providers/News"
	"context"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		baseProvider News.IProviderNews
		viper        *viper.Viper
		news         *Adapters.NewsAdapter
	}
	tests := []struct {
		name string
		args args
		want *NewsApiProvider
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProvider(tt.args.baseProvider, tt.args.viper, tt.args.news); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsApiProvider_ResponseReaderFunc(t *testing.T) {
	type fields struct {
		IProviderNews News.IProviderNews
		Viper         *viper.Viper
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
				Config:        tt.fields.Viper,
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
		Viper         *viper.Viper
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
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &NewsApiProvider{
				IProviderNews: tt.fields.IProviderNews,
				Config:        tt.fields.Viper,
				Adapter:       tt.fields.Adapter,
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
