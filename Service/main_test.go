package main

import (
	"Providers/Configs"
	"context"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *viper.Viper
		want1   *Configs.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoadConfig() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_handleRequest(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleRequest(tt.args.ctx)
		})
	}
}

func Test_loadProviderConfig(t *testing.T) {
	type args struct {
		v   *viper.Viper
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    Configs.ProviderConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadProviderConfig(tt.args.v, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadProviderConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadProviderConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
