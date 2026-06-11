package usecase

import (
	"context"
	"time"

	"github.com/gokhan/orderly/internal/domain"
	"github.com/gokhan/orderly/internal/repository/db"
	"github.com/gokhan/orderly/pkg/config"
	"github.com/gokhan/orderly/pkg/utils"
)

type userUseCase struct {
	store  db.Store
	config config.Config
}

func NewUserUseCase(store db.Store, config config.Config) domain.UserUseCase {
	return &userUseCase{
		store:  store,
		config: config,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (domain.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return domain.UserResponse{}, err
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := u.store.CreateUser(ctx, arg)
	if err != nil {
		return domain.UserResponse{}, err
	}

	return domain.NewUserResponse(user), nil
}

func (u *userUseCase) LoginUser(ctx context.Context, req domain.LoginUserRequest) (domain.LoginUserResponse, error) {
	user, err := u.store.GetUser(ctx, req.Username)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	duration, err := time.ParseDuration(u.config.AccessTokenDuration)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	accessToken, err := utils.CreateToken(user.Username, duration, u.config.TokenSymmetricKey)
	if err != nil {
		return domain.LoginUserResponse{}, err
	}

	resp := domain.LoginUserResponse{
		AccessToken: accessToken,
		User:        domain.NewUserResponse(user),
	}

	return resp, nil
}
