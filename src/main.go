package main

import (
	"context"
	"database/sql"
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
	setupLog  = ctrl.Log.WithName("setup")
	log       = logf.Log.WithName("user-service")
	port      = "3550"
	loginType = []string{"discord", "telegram", "twitter"}
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
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, "DB connection Error")
	}
	defer db.Close()

	var id int
	var result sql.Result

	// TODO use MySQL's unique function
	query := fmt.Sprintf("SELECT id FROM user WHERE %s = '%s'", loginType[request.LoginType], request.UserID)
	err = db.QueryRowContext(ctx, query).Scan(&id)
	if err == sql.ErrNoRows {
		query = fmt.Sprintf("INSERT INTO user(type, %s) VALUES(%d, '%s')", loginType[request.LoginType], request.LoginType, request.UserID)
		result, err = db.ExecContext(ctx, query)
		if err != nil {
			log.Error(err, "create user error")
			return nil, status.Errorf(codes.Internal, "create user error")
		}
		n, err := result.RowsAffected()
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if n == 1 {
			response := &pb.GetUserResponse{}
			query = fmt.Sprintf("SELECT * FROM user WHERE %s = '%s'", loginType[request.LoginType], request.UserID)
			if err = db.QueryRowContext(ctx, query).Scan(&response.UserID, &response.LoginType, &response.TelegramID, &response.DiscordID, &response.TwitterID); err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			return response, nil
		}
	}

	if err == nil {
		log.Info("exising id error")
		return nil, status.Errorf(codes.InvalidArgument, "already registered id")
	}
	log.Error(err, "")
	return nil, status.Errorf(codes.Internal, err.Error())
}

func (u *userService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.GetUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, "DB connection Error")
	}
	defer db.Close()

	if request.LoginType > 2 {
		err = fmt.Errorf("wrong login type")
		log.Error(err, "")
		return nil, status.Errorf(codes.InvalidArgument, "wrong login type")
	}

	query := fmt.Sprintf("SELECT * FROM user WHERE type = %d AND %s = '%s'", request.LoginType, loginType[request.LoginType], request.UserID)
	row := db.QueryRowContext(ctx, query)

	response := &pb.GetUserResponse{}

	err = row.Scan(&response.UserID, &response.LoginType, &response.TelegramID, &response.DiscordID, &response.TwitterID)

	if err == sql.ErrNoRows {
		return u.CreateUser(ctx, &pb.CreateUserRequest{
			UserID:    request.GetUserID(),
			LoginType: request.GetLoginType(),
		})
	}
	if err != nil {
		log.Error(err, "")
		return nil, status.Errorf(codes.InvalidArgument, "scan error")
	}
	// TODO implement login
	return response, nil
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
