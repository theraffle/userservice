package userservice

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/theraffle/userservice/src/database"
	"github.com/theraffle/userservice/src/genproto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	log       = logf.Log.WithName("user-service")
	loginType = []string{"discord", "telegram", "twitter"}
)

// UserService is a struct that have user service gRPC methods
type UserService struct{}

// CreateUser creates a user in user table and returns generated user ID
func (u *UserService) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.LoginUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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

// GetUser returns user's info from user table
func (u *UserService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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

// LoginUser returns current user's ID
func (u *UserService) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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

// LogoutUser exists for just in case
func (u *UserService) LogoutUser(_ context.Context, _ *pb.LogoutUserRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// UpdateUser updates info of requested user and returns updated info
func (u *UserService) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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

// CreateUserWallet creates new user's wallet in user_wallet table
func (u *UserService) CreateUserWallet(ctx context.Context, request *pb.CreateUserWalletRequest) (*pb.Empty, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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

// GetUserWallet returns a user's entire wallets
func (u *UserService) GetUserWallet(ctx context.Context, request *pb.GetUserWalletRequest) (*pb.GetUserWalletResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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

// CreateUserProject creates a user's project in user_project table
func (u *UserService) CreateUserProject(ctx context.Context, request *pb.CreateUserProjectRequest) (*pb.Empty, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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

// GetUserProject returns user's entire related projects
func (u *UserService) GetUserProject(ctx context.Context, request *pb.GetUserProjectRequest) (*pb.GetUserProjectResponse, error) {
	db, err := database.Connect()
	if err != nil {
		log.Error(err, "database connection error")
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err, "db closing error")
			os.Exit(1)
		}
	}()

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
