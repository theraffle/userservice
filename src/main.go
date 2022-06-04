package main

import (
	"context"
	proto "github.com/theraffle/backend/src/genproto/pb"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	proto.RegisterUserManagerServer(grpcServer, &userManagerServer{})
}

type userManagerServer struct{}

func (u *userManagerServer) LogoutUser(ctx context.Context, request *proto.LogoutUserRequest) (*proto.Error, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userManagerServer) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.Error, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userManagerServer) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userManagerServer) LoginUser(ctx context.Context, request *proto.LoginUserRequest) (*proto.Error, error) {
	//TODO implement me
	panic("implement me")
}
