package main

import (
	"context"
	"testing"
)

func Test_handleRequest(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_handleRequest_success",
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleRequest(tt.args.ctx)
		})
	}
}
