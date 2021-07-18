package main

import (
	"context"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	"reflect"
	"testing"
	"time"
)

func TestServer_Infer(t *testing.T) {
	type args struct {
		ctx context.Context
		req *deepthought.InferRequest
	}
	var (
		tests = []struct {
			name    string
			args    args
			want    *deepthought.InferResponse
			wantErr bool
		}{
			{
				name: "Life",
				args: args{
					ctx: context.TODO(),
					req: &deepthought.InferRequest{
						Query: "Life",
					},
				},
				want:    &deepthought.InferResponse{Answer: 42},
				wantErr: false,
			},
			{
				name: "Universe",
				args: args{
					ctx: context.TODO(),
					req: &deepthought.InferRequest{
						Query: "Universe",
					},
				},
				want:    &deepthought.InferResponse{Answer: 42},
				wantErr: false,
			},
			{
				name: "Everything",
				args: args{
					ctx: context.TODO(),
					req: &deepthought.InferRequest{
						Query: "Everything",
					},
				},
				want:    &deepthought.InferResponse{Answer: 42},
				wantErr: false,
			},
			{
				name: "other",
				args: args{
					ctx: context.TODO(),
					req: &deepthought.InferRequest{
						Query: "hogehoge",
					},
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "deadline",
				args: args{
					ctx: func() context.Context {
						ctx, _ := context.WithTimeout(context.TODO(), 1*time.Millisecond)
						return ctx
					}(),
					req: &deepthought.InferRequest{
						Query: "Life",
					},
				},
				want:    nil,
				wantErr: true,
			},
		}
	)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &DeepthoughtServer{}
			got, err := s.Infer(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Infer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Infer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
