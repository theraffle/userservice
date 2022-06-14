/*
 Copyright 2022 The Raffle Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

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

type userService struct{}

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

// TODO id 중복체크

func (u *userService) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.LoginUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO user(type, %s) VALUES(%d, '%s')", loginType[request.LoginType], request.LoginType, request.UserID)
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Error(err, "create user error")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	n, _ := result.RowsAffected()
	if n == 1 {
		response := &pb.LoginUserResponse{}
		query = fmt.Sprintf("SELECT id FROM user WHERE %s = '%s'", loginType[request.LoginType], request.UserID)
		if err = db.QueryRowContext(ctx, query).Scan(&response.UserID); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return response, nil
	}

	log.Error(err, "check database user table")
	return nil, status.Errorf(codes.Internal, "user insertion error: check database user table")
}

func (u *userService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM user WHERE id = %d", request.UserID)
	row := db.QueryRowContext(ctx, query)

	response := &pb.GetUserResponse{}
	err = row.Scan(&response.UserID, &response.LoginType, &response.TelegramID, &response.DiscordID, &response.TwitterID)

	if err != nil {
		log.Error(err, "getting user from db error")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return response, nil
}

func (u *userService) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	if request.LoginType > 2 {
		err = fmt.Errorf("wrong login type")
		log.Error(err, "")
		return nil, status.Errorf(codes.InvalidArgument, "wrong login type")
	}

	query := fmt.Sprintf("SELECT id FROM user WHERE type = %d AND %s = '%s'", request.LoginType, loginType[request.LoginType], request.UserID)
	row := db.QueryRowContext(ctx, query)

	response := &pb.LoginUserResponse{}

	err = row.Scan(&response.UserID)

	if err == sql.ErrNoRows {
		return u.CreateUser(ctx, &pb.CreateUserRequest{
			UserID:    request.GetUserID(),
			LoginType: request.GetLoginType(),
		})
	}
	if err != nil {
		log.Error(err, "")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// TODO implement login
	return response, nil
}

func (u *userService) LogoutUser(_ context.Context, _ *pb.LogoutUserRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (u *userService) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE user SET telegram = '%s', discord = '%s', twitter = '%s' WHERE id = %d", request.TelegramID, request.DiscordID, request.TwitterID, request.UserID)
	result, err := db.ExecContext(ctx, query)

	if err != nil {
		log.Error(err, "update user error")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	n, _ := result.RowsAffected()
	if n == 1 {
		response := &pb.GetUserResponse{}
		query = fmt.Sprintf("SELECT * FROM user WHERE id = '%d'", request.UserID)
		if err = db.QueryRowContext(ctx, query).Scan(&response.UserID, &response.LoginType, &response.TelegramID, &response.DiscordID, &response.TwitterID); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return response, nil
	}

	log.Error(err, "check database user table")
	return nil, status.Errorf(codes.Internal, "user update error: check database user table")
}

func (u *userService) CreateUserWallet(ctx context.Context, request *pb.CreateUserWalletRequest) (*pb.Empty, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	_, err = u.GetUser(ctx, &pb.GetUserRequest{UserID: request.Wallet.UserID})
	if err != nil {
		log.Error(err, fmt.Sprintf("cannot find user id %d", request.Wallet.UserID))
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	query := fmt.Sprintf("INSERT INTO user_wallet(user_id, chain_id, address) VALUES(%d, %d, '%s')", request.Wallet.UserID, request.Wallet.ChainID, request.Wallet.Address)
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Error(err, "create user wallet error")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	n, _ := result.RowsAffected()
	if n == 1 {
		return &pb.Empty{}, nil
	}

	log.Error(err, "check database user_wallet table")
	return nil, status.Errorf(codes.Internal, "user_wallet insertion error: check database user_wallet table")
}

func (u *userService) GetUserWallet(ctx context.Context, request *pb.GetUserWalletRequest) (*pb.GetUserWalletResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM user_wallet WHERE user_id = %d", request.UserID)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Error(err, "get wallet query error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	var wallets []*pb.UserWallet
	for rows.Next() {
		var wallet *pb.UserWallet
		err = rows.Scan(&wallet.UserID, &wallet.ChainID, &wallet.Address)
		if err != nil {
			log.Error(err, "get user wallet scan error")
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		wallets = append(wallets, wallet)
	}
	return &pb.GetUserWalletResponse{Wallets: wallets}, nil
}

func (u *userService) CreateUserProject(ctx context.Context, request *pb.CreateUserProjectRequest) (*pb.Empty, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	_, err = u.GetUser(ctx, &pb.GetUserRequest{UserID: request.UserID})
	if err != nil {
		log.Error(err, fmt.Sprintf("cannot find user id %d", request.UserID))
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	query := fmt.Sprintf("INSERT INTO user_project(user_id, prject_id, chain_id, address) VALUES(%d, %d, %d, '%s')", request.UserID, request.ProjectID, request.ChainID, request.Address)
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Error(err, "create user project error")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	n, _ := result.RowsAffected()
	if n == 1 {
		return &pb.Empty{}, nil
	}

	log.Error(err, "check database user_project table")
	return nil, status.Errorf(codes.Internal, "user_project insertion error: check database user_project table")
}

func (u *userService) GetUserProject(ctx context.Context, request *pb.GetUserProjectRequest) (*pb.GetUserProjectResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM user_project WHERE user_id = %d", request.UserID)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Error(err, "get project query error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	var projects []int64
	for rows.Next() {
		var project int64
		err = rows.Scan(&project)
		if err != nil {
			log.Error(err, "get user project scan error")
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		projects = append(projects, project)
	}
	return &pb.GetUserProjectResponse{Projects: projects}, nil
}
