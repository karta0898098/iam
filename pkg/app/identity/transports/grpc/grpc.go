package grpc

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	pb "github.com/karta0898098/iam/pb/identity"
	"github.com/karta0898098/iam/pkg/app/identity/endpoints"
	"github.com/karta0898098/iam/pkg/app/identity/entity"
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

	return &endpoints.SigninRequest{
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

// encodeGRPCSigninRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Signin request to a gRPC Signin request. Primarily useful in a client.
func encodeGRPCSigninRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*endpoints.SigninRequest)
	return &pb.SigninReq{
		Username:  req.Username,
		Password:  req.Password,
		IPAddress: req.IPAddress,
		Device: &pb.Device{
			Model:     req.Device.Model,
			Name:      req.Device.Name,
			OSVersion: req.Device.OSVersion,
		},
		IdpProvider: req.IdpProvider,
	}, nil
}

func decodeGRPCSigninResponse(ctx context.Context, grpcReply interface{}) (response interface{}, err error) {
	reply := grpcReply.(*pb.SigninResp)
	return &endpoints.SigninResponse{
		IDToken:      reply.IDToken,
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
	}, nil
}

// decodeGRPCSignupRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCSignupRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SignupReq)

	return &endpoints.SignupRequest{
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

// encodeGRPCSignupResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCSignupResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(*endpoints.SignupResponse)
	return &pb.SignupResp{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		IDToken:      reply.IDToken,
	}, nil
}

// encodeGRPCSignupRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Signup request to a gRPC Signup request. Primarily useful in a client.
func encodeGRPCSignupRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*endpoints.SigninRequest)
	return &pb.SigninReq{
		Username:  req.Username,
		Password:  req.Password,
		IPAddress: req.IPAddress,
		Device: &pb.Device{
			Model:     req.Device.Model,
			Name:      req.Device.Name,
			OSVersion: req.Device.OSVersion,
		},
		IdpProvider: req.IdpProvider,
	}, nil
}

// decodeGRPCSignupResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Signup reply to a user-domain Sum response. Primarily useful in a client.
func decodeGRPCSignupResponse(ctx context.Context, grpcReply interface{}) (response interface{}, err error) {
	reply := grpcReply.(*pb.SignupResp)
	return &endpoints.SignupResponse{
		IDToken:      reply.IDToken,
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
	}, nil
}
