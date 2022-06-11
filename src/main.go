package main

import (
	"context"
	"github.com/theraffle/backend/src/genproto/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	run()
}

func run() string {
	l, err := net.Listen("tcp", "3550")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()

	svc := &userService{}

	pb.RegisterUserServiceServer(srv, svc)
	go srv.Serve(l)
	return l.Addr().String()
}

type userService struct{}

func (u *userService) LogoutUser(ctx context.Context, request *pb.LogoutUserRequest) (*pb.Error, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.Error, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.Error, error) {
	//TODO implement me
	panic("implement me")
}
