package main

import (
	"context"
	"github.com/care0717/deepthought-api/grpc/proto/auth"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	repository2 "github.com/care0717/deepthought-api/grpc/server/repository"
	service2 "github.com/care0717/deepthought-api/grpc/server/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type DeepthoughtServer struct {
	deepthought.UnimplementedComputeServer
}

var _ deepthought.ComputeServer = &DeepthoughtServer{}

func (s *DeepthoughtServer) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	if req.Silent {
		return nil
	}
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(1 * time.Second):
		}
		if err := stream.Send(&deepthought.BootResponse{
			Message: "I think therefore i am.",
			Ts:      timestamppb.Now(),
		}); err != nil {
			return err
		}
	}
}

func (s *DeepthoughtServer) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	switch req.Query {
	case "Life", "Universe", "Everything":
	default:
		return nil, status.Error(codes.InvalidArgument, "Contemplate your query")
	}

	deadline, ok := ctx.Deadline()

	if !ok || time.Until(deadline) > 750*time.Millisecond {
		time.Sleep(750 * time.Millisecond)
		return &deepthought.InferResponse{
			Answer: 42,
		}, nil
	}

	return nil, status.Error(codes.DeadlineExceeded, "It would take longer")
}

type AuthServer struct {
	userStore  repository2.UserStore
	jwtManager *service2.JWTManager
	auth.UnimplementedAuthServer
}

var _ auth.AuthServer = &AuthServer{}

func NewAuthServer(userStore repository2.UserStore, jwtManager *service2.JWTManager) *AuthServer {
	return &AuthServer{userStore: userStore, jwtManager: jwtManager}
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := s.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &auth.LoginResponse{AccessToken: token}
	return res, nil
}
