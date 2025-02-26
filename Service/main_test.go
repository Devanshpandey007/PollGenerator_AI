package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"reflect"
	"testing"
)

func Test_handleRequest(t *testing.T) {
	type args struct {
		ctx context.Context
		in1 events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Test_handleRequest_success",
			args: args{
				ctx: context.Background(),
				in1: events.APIGatewayProxyRequest{},
			},
			want: events.APIGatewayProxyResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleRequest(tt.args.ctx, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
