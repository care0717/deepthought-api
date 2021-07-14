package main

import (
	"context"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Server struct {
	deepthought.UnimplementedComputeServer
}

var _ deepthought.ComputeServer = &Server{}

func (s *Server) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(1 * time.Second):
		}
		if err := stream.Send(&deepthought.BootResponse{Message: "I think therefore i am."}); err != nil {
			return err
		}
	}
}

func (s *Server) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	switch req.Query {
	case "Life", "Universe", "Everything":
	default:
		return nil, status.Error(codes.InvalidArgument, "Contemplate your query")
	}

	deadline, ok := ctx.Deadline()

	if !ok || time.Until(deadline) > 750*time.Millisecond {
		time.Sleep(750 * time.Millisecond)
		return &deepthought.InferResponse{
			Answer:      42,
			Description: []string{"I checked it"},
		}, nil
	}

	return nil, status.Error(codes.DeadlineExceeded, "It would take longer")
}
