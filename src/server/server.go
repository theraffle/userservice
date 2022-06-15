package server

import (
	"fmt"
	"github.com/theraffle/userservice/src/genproto/pb"
	"github.com/theraffle/userservice/src/server/userservice"
	"google.golang.org/grpc"
	"net"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	log = logf.Log.WithName("server")
)

// Server is an interface of server
type Server interface {
	Start(string)
}

type userServiceServer struct {
	grpcServer  *grpc.Server
	userService *userservice.UserService
}

// New returns new user service gRPC server
func New() Server {
	server := new(userServiceServer)
	server.grpcServer = grpc.NewServer()
	server.userService = &userservice.UserService{}

	return server
}

// Start starts gRPC server on input port
func (s *userServiceServer) Start(port string) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Error(err, "cannot launch gRPC server")
		os.Exit(1)
	}

	pb.RegisterUserServiceServer(s.grpcServer, s.userService)
	if err := s.grpcServer.Serve(l); err != nil {
		log.Error(err, "cannot launch gRPC server")
		os.Exit(1)
	}
}
