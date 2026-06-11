package usecase

import (
	"context"
	"testing"

	"github.com/gokhan/orderly/internal/domain"
	"github.com/gokhan/orderly/internal/repository/db"
	mockdb "github.com/gokhan/orderly/internal/repository/db/mock"
	"github.com/gokhan/orderly/pkg/config"
	"github.com/gokhan/orderly/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	user := db.User{
		Username: "testuser",
		FullName: "Test User",
		Email:    "test@example.com",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// Build stub
	store.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Times(1).
		Return(user, nil)

	uc := NewUserUseCase(store, config.Config{})
	res, err := uc.CreateUser(context.Background(), domain.CreateUserRequest{
		Username: user.Username,
		Password: "password123",
		FullName: user.FullName,
		Email:    user.Email,
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, res.Username)
}

func TestLoginUser(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")
	user := db.User{
		ID:             1,
		Username:       "testuser",
		HashedPassword: hashedPassword,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().
		GetUser(gomock.Any(), gomock.Eq(user.Username)).
		Times(1).
		Return(user, nil)

	cfg := config.Config{
		TokenSymmetricKey:   "v2.local.6v2Z5u6B6n8m6B6v2Z5u6B6n8m6B6v2Z",
		AccessTokenDuration: "15m",
	}

	uc := NewUserUseCase(store, cfg)
	res, err := uc.LoginUser(context.Background(), domain.LoginUserRequest{
		Username: user.Username,
		Password: "password123",
	})

	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
	require.Equal(t, user.Username, res.User.Username)
}
