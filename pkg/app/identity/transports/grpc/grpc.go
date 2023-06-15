package grpc

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/karta0898098/iam/pkg/app/identity/endpoints"
	"github.com/karta0898098/iam/pkg/app/identity/entity"

	pb "github.com/karta0898098/iam/pb/identity"
)

type grpcServer struct {
	signin grpctransport.Handler
	signup grpctransport.Handler
}

func (g *grpcServer) Signin(ctx context.Context, req *pb.SigninReq) (*pb.SigninResp, error) {
	_, rp, err := g.signin.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	reply := (rp).(*pb.SigninResp)

	return reply, nil
}

func (g *grpcServer) Signup(ctx context.Context, req *pb.SignupReq) (*pb.SignupResp, error) {

	_, rp, err := g.signup.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	reply := (rp).(*pb.SignupResp)
	return reply, nil
}

func MakeGRPCServer(endpoints endpoints.Endpoints) (req pb.IdentityServiceServer) {
	options := []grpctransport.ServerOption{}

	return &grpcServer{
		signin: grpctransport.NewServer(
			endpoints.SigninEndpoint,
			decodeGRPCSigninRequest,
			encodeGRPCSigninResponse,
			options...,
		),
		signup: grpctransport.NewServer(
			endpoints.SignupEndpoint,
			decodeGRPCSignupRequest,
			encodeGRPCSignupResponse,
			options...,
		),
	}
}

// decodeGRPCSigninRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCSigninRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SigninReq)

	return endpoints.SigninRequest{
		Username:  req.Username,
		Password:  req.Password,
		IPAddress: req.IPAddress,
		Device: entity.Device{
			Model:     req.Device.Model,
			Name:      req.Device.Name,
			OSVersion: req.Device.OSVersion,
		},
	}, nil
}

// encodeGRPCSigninResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCSigninResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(*endpoints.SigninResponse)
	return &pb.SigninResp{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		IDToken:      reply.IDToken,
	}, nil
}

// decodeGRPCTicRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCSignupRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SignupReq)

	return endpoints.SignupRequest{
		Username:  req.Username,
		Password:  req.Password,
		Nickname:  req.Nickname,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Platform:  req.Platform,
		IPAddress: req.IPAddress,
		Device: entity.Device{
			Model:     req.Device.Model,
			Name:      req.Device.Name,
			OSVersion: req.Device.OSVersion,
		},
	}, nil
}

// encodeGRPCTicResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCSignupResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(*endpoints.SignupResponse)
	return &pb.SignupResp{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		IDToken:      reply.IDToken,
	}, nil
}