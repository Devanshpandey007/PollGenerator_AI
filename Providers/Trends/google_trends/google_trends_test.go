package google_trends

import (
	Adapters "Providers/Adapters/Trends"
	"Providers/Configs"
	"Providers/Providers/Trends"
	"context"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestGoogleTrendsProvider_GetConfig(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Viper           *viper.Viper
		Adapter         *Adapters.TrendsAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   Configs.ProviderConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Viper,
				Adapter:         tt.fields.Adapter,
			}
			if got := p.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleTrendsProvider_GetProviderName(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Viper           *viper.Viper
		Adapter         *Adapters.TrendsAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Viper,
				Adapter:         tt.fields.Adapter,
			}
			if got := p.GetProviderName(); got != tt.want {
				t.Errorf("GetProviderName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleTrendsProvider_GetTrend(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Viper           *viper.Viper
		Adapter         *Adapters.TrendsAdapter
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Viper,
				Adapter:         tt.fields.Adapter,
			}
			got, err := p.GetTrend(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTrend() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleTrendsProvider_GetTrendFromRegion(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Viper           *viper.Viper
		Adapter         *Adapters.TrendsAdapter
	}
	type args struct {
		ctx    context.Context
		region string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Viper,
				Adapter:         tt.fields.Adapter,
			}
			got, err := p.GetTrendFromRegion(tt.args.ctx, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrendFromRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTrendFromRegion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleTrendsProvider_GetTrends(t *testing.T) {
	type fields struct {
		IProviderTrends Trends.IProviderTrends
		Viper           *viper.Viper
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GoogleTrendsProvider{
				IProviderTrends: tt.fields.IProviderTrends,
				Config:          tt.fields.Viper,
				Adapter:         tt.fields.Adapter,
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
		Viper           *viper.Viper
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
				Config:          tt.fields.Viper,
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
		Viper           *viper.Viper
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
				Config:          tt.fields.Viper,
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
		Viper           *viper.Viper
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
				Config:          tt.fields.Viper,
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
		baseProvider Trends.IProviderTrends
		viper        *viper.Viper
		trends       *Adapters.TrendsAdapter
	}
	tests := []struct {
		name string
		args args
		want *GoogleTrendsProvider
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProvider(tt.args.baseProvider, tt.args.viper, tt.args.trends); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
