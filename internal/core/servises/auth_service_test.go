package servises

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/domain"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/ports/dto"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/ports/mocks"
	"github.com/ttodoshi/code-typing-auth-service/pkg/jwt"
	"github.com/ttodoshi/code-typing-auth-service/pkg/logging/nop"
	. "github.com/ttodoshi/code-typing-auth-service/pkg/password"
	"os"
	"testing"
)

func TestRegister(t *testing.T) {
	var log = nop.GetLogger()
	var err error
	err = os.Setenv("ACCESS_TOKEN_EXP", "300")
	err = os.Setenv("REFRESH_TOKEN_EXP", "1209600")
	// mocks
	userRepo := new(mocks.UserRepository)
	tokenRepo := new(mocks.RefreshTokenRepository)
	eventDispatcher := new(mocks.EventDispatcher)

	userRepo.
		On("GetUserByNickname", "already_taken").
		Return(domain.User{}, nil)
	userRepo.
		On("GetUserByNickname", mock.Anything).
		Return(domain.User{}, fmt.Errorf(""))
	userRepo.
		On("GetUserByEmail", "already_taken").
		Return(domain.User{}, nil)
	userRepo.
		On("GetUserByEmail", mock.Anything).
		Return(domain.User{}, fmt.Errorf(""))
	userRepo.
		On("SaveUser", mock.Anything).
		Return(domain.User{}, nil)
	tokenRepo.
		On("CreateRefreshToken", mock.Anything).
		Return(gofakeit.UUID(), nil)
	eventDispatcher.
		On(
			"Dispatch",
			mock.Anything,
		).Return()

	// service
	authService := NewAuthService(userRepo, tokenRepo, eventDispatcher, log)

	t.Run("successful registration", func(t *testing.T) {
		_, _, err = authService.Register(dto.RegisterRequestDto{
			Nickname: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
		}, gofakeit.UUID())
		assert.NoError(t, err)
	})
	t.Run("unsuccessful registration due to nickname already taken", func(t *testing.T) {
		_, _, err = authService.Register(dto.RegisterRequestDto{
			Nickname: "already_taken",
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 4),
		}, gofakeit.UUID())
		assert.Error(t, err)
	})
	t.Run("unsuccessful registration due to email already taken", func(t *testing.T) {
		_, _, err = authService.Register(dto.RegisterRequestDto{
			Nickname: gofakeit.Username(),
			Email:    "already_taken",
			Password: gofakeit.Password(true, true, true, true, false, 4),
		}, gofakeit.UUID())
		assert.Error(t, err)
	})
	userRepo.AssertExpectations(t)
	eventDispatcher.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	var log = nop.GetLogger()
	var err error
	err = os.Setenv("ACCESS_TOKEN_EXP", "300")
	err = os.Setenv("REFRESH_TOKEN_EXP", "1209600")
	// mocks
	userRepo := new(mocks.UserRepository)
	tokenRepo := new(mocks.RefreshTokenRepository)
	eventDispatcher := new(mocks.EventDispatcher)

	password := gofakeit.Password(true, true, true, true, false, 8)
	hashPassword, err := HashPassword(password)
	user := domain.User{
		Nickname: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: hashPassword,
	}
	userRepo.
		On("GetUserByNickname", user.Nickname).
		Return(user, nil)
	userRepo.
		On("GetUserByEmail", user.Email).
		Return(user, nil)
	userRepo.
		On("GetUserByNickname", mock.AnythingOfType("string")).
		Return(domain.User{}, fmt.Errorf(""))
	userRepo.
		On("GetUserByEmail", mock.AnythingOfType("string")).
		Return(domain.User{}, fmt.Errorf(""))

	tokenRepo.
		On("CreateRefreshToken", mock.Anything).
		Return(gofakeit.UUID(), nil)
	eventDispatcher.
		On(
			"Dispatch",
			mock.Anything,
		).Return()

	// service
	authService := NewAuthService(userRepo, tokenRepo, eventDispatcher, log)

	t.Run("successful login by nickname", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    user.Nickname,
			Password: password,
		}, gofakeit.UUID())
		assert.NoError(t, err)
	})
	t.Run("successful login by email", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    user.Email,
			Password: password,
		}, gofakeit.UUID())
		assert.NoError(t, err)
	})
	t.Run("unsuccessful login due to invalid email", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    "invalid_email",
			Password: password,
		}, gofakeit.UUID())
		assert.Error(t, err)
	})
	t.Run("unsuccessful login due to invalid nickname", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    "invalid_nickname",
			Password: password,
		}, gofakeit.UUID())
		assert.Error(t, err)
	})
	t.Run("unsuccessful login due to invalid password", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    user.Nickname,
			Password: "invalid_password",
		}, gofakeit.UUID())
		assert.Error(t, err)
	})
	userRepo.AssertExpectations(t)
	eventDispatcher.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRefresh(t *testing.T) {
	var log = nop.GetLogger()
	var err error
	err = os.Setenv("ACCESS_TOKEN_EXP", "300")
	err = os.Setenv("REFRESH_TOKEN_EXP", "1209600")
	// mocks
	userRepo := new(mocks.UserRepository)
	tokenRepo := new(mocks.RefreshTokenRepository)
	eventDispatcher := new(mocks.EventDispatcher)

	password := gofakeit.Password(true, true, true, true, false, 8)
	hashPassword, err := HashPassword(password)
	user := domain.User{
		Nickname: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: hashPassword,
	}

	refresh, err := jwt.GenerateRefreshJWT(user.ID.Hex())
	tokenRepo.
		On("GetRefreshToken", refresh).
		Return(domain.RefreshToken{
			User:  user.ID,
			Token: refresh,
		}, nil)
	tokenRepo.
		On("GetRefreshToken", mock.AnythingOfType("string")).
		Return(domain.RefreshToken{}, fmt.Errorf(""))
	tokenRepo.
		On("UpdateRefreshToken", refresh, mock.Anything).
		Return(domain.RefreshToken{}, nil)

	userRepo.
		On("GetUserByID", user.ID.Hex()).
		Return(user, nil)

	// service
	authService := NewAuthService(userRepo, tokenRepo, eventDispatcher, log)

	t.Run("successful refresh", func(t *testing.T) {
		_, _, err = authService.Refresh(refresh)
		assert.NoError(t, err)
	})
	t.Run("unsuccessful refresh due to invalid refresh token", func(t *testing.T) {
		_, _, err = authService.Refresh("invalid_refresh_token")
		assert.Error(t, err)
	})
	userRepo.AssertExpectations(t)
	eventDispatcher.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestLogout(t *testing.T) {
	var log = nop.GetLogger()
	// mocks
	userRepo := new(mocks.UserRepository)
	tokenRepo := new(mocks.RefreshTokenRepository)
	eventDispatcher := new(mocks.EventDispatcher)

	refresh, err := jwt.GenerateRefreshJWT(gofakeit.UUID())
	tokenRepo.
		On("DeleteRefreshToken", refresh).
		Return(nil)
	tokenRepo.
		On("DeleteRefreshToken", mock.AnythingOfType("string")).
		Return(fmt.Errorf(""))

	// service
	authService := NewAuthService(userRepo, tokenRepo, eventDispatcher, log)

	t.Run("successful logout", func(t *testing.T) {
		authService.Logout(refresh)
		assert.NoError(t, err)
	})
	userRepo.AssertExpectations(t)
	eventDispatcher.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}
