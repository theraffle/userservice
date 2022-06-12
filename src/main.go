package main

import (
	"context"
	"fmt"
	"github.com/theraffle/backend/src/genproto/pb"
	"github.com/theraffle/backend/src/internal/database"
	"github.com/theraffle/backend/src/internal/logrotate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	setupLog = ctrl.Log.WithName("setup")
	log      = logf.Log.WithName("user-service")
	port     = "3550"
)

func main() {
	// Set log rotation
	logFile, err := logrotate.LogFile()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() {
		_ = logFile.Close()
	}()
	logWriter := io.MultiWriter(logFile, os.Stdout)
	ctrl.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(logWriter)))
	if err := logrotate.StartRotate("0 0 1 * * ?"); err != nil {
		setupLog.Error(err, "")
		os.Exit(1)
	}
	// Set ports
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	run(port)
}

func run(port string) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Error(err, "cannot launch gRPC server")
		os.Exit(1)
	}
	srv := grpc.NewServer()

	svc := &userService{}

	pb.RegisterUserServiceServer(srv, svc)
	if err := srv.Serve(l); err != nil {
		log.Error(err, "cannot launch gRPC server")
		os.Exit(1)
	}
}

type userService struct{}

func (u *userService) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	log.Info("test")
	_, err := database.Connect()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "DB connection Error")
	}
	return nil, nil
}

func (u *userService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) LogoutUser(ctx context.Context, request *pb.LogoutUserRequest) (*pb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) GetUserWallet(ctx context.Context, request *pb.GetUserWalletRequest) (*pb.GetUserWalletResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) GetUserProject(ctx context.Context, response *pb.GetUserProjectResponse) (*pb.GetUserProjectResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) UpdateUserWallet(ctx context.Context, request *pb.UpdateUserWalletRequest) (*pb.GetUserWalletRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) UpdateUserProject(ctx context.Context, request *pb.UpdateUserProjectRequest) (*pb.GetUserProjectResponse, error) {
	//TODO implement me
	panic("implement me")
}
