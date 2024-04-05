package servises

import (
	"code-typing-auth-service/internal/adapters/dto"
	"code-typing-auth-service/internal/core/domain"
	"code-typing-auth-service/internal/core/errors"
	"code-typing-auth-service/internal/core/ports/mocks"
	"code-typing-auth-service/internal/core/utils"
	"code-typing-auth-service/pkg/logging/nop"
	. "code-typing-auth-service/pkg/password"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	resultsMigrator := new(mocks.ResultsMigrator)

	userRepo.
		On("SaveUser", mock.MatchedBy(func(arg interface{}) bool {
			u := arg.(domain.User)
			return u.Nickname == "already_taken" || u.Email == "already_taken"
		})).
		Maybe().
		Return(domain.User{}, &errors.AlreadyExistsError{})
	userRepo.
		On("SaveUser", mock.Anything).
		Maybe().
		Return(domain.User{}, nil)
	tokenRepo.
		On("SaveRefreshToken", mock.Anything).
		Return(gofakeit.UUID(), nil)
	resultsMigrator.
		On(
			"MigrateSessionResults",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return()

	// service
	authService := NewAuthService(userRepo, tokenRepo, resultsMigrator, log)

	t.Run("successful registration", func(t *testing.T) {
		_, _, err = authService.Register(dto.RegisterRequestDto{
			Nickname: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
		}, gofakeit.UUID())
		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
	t.Run("unsuccessful registration due to nickname already taken", func(t *testing.T) {
		_, _, err = authService.Register(dto.RegisterRequestDto{
			Nickname: "already_taken",
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 4),
		}, gofakeit.UUID())
		assert.Error(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
	t.Run("unsuccessful registration due to email already taken", func(t *testing.T) {
		_, _, err = authService.Register(dto.RegisterRequestDto{
			Nickname: gofakeit.Username(),
			Email:    "already_taken",
			Password: gofakeit.Password(true, true, true, true, false, 4),
		}, gofakeit.UUID())
		assert.Error(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	var log = nop.GetLogger()
	var err error
	err = os.Setenv("ACCESS_TOKEN_EXP", "300")
	err = os.Setenv("REFRESH_TOKEN_EXP", "1209600")
	// mocks
	userRepo := new(mocks.UserRepository)
	tokenRepo := new(mocks.RefreshTokenRepository)
	resultsMigrator := new(mocks.ResultsMigrator)

	password := gofakeit.Password(true, true, true, true, false, 8)
	hashPassword, err := HashPassword(password)
	user := domain.User{
		Nickname: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: hashPassword,
	}
	userRepo.
		On("GetUserByNickname", user.Nickname).
		Maybe().
		Return(user, nil)
	userRepo.
		On("GetUserByEmail", user.Email).
		Maybe().
		Return(user, nil)
	userRepo.
		On("GetUserByNickname", mock.AnythingOfType("string")).
		Maybe().
		Return(domain.User{}, &errors.NotFoundError{})
	userRepo.
		On("GetUserByEmail", mock.AnythingOfType("string")).
		Maybe().
		Return(domain.User{}, &errors.NotFoundError{})

	tokenRepo.
		On("SaveRefreshToken", mock.Anything).
		Return(gofakeit.UUID(), nil)
	resultsMigrator.
		On(
			"MigrateSessionResults",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return()

	// service
	authService := NewAuthService(userRepo, tokenRepo, resultsMigrator, log)

	t.Run("successful login by nickname", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    user.Nickname,
			Password: password,
		}, gofakeit.UUID())
		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
	t.Run("successful login by email", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    user.Email,
			Password: password,
		}, gofakeit.UUID())
		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
	t.Run("unsuccessful login due to invalid email", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    "invalid_email",
			Password: password,
		}, gofakeit.UUID())
		assert.Error(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
	t.Run("unsuccessful login due to invalid nickname", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    "invalid_nickname",
			Password: password,
		}, gofakeit.UUID())
		assert.Error(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
	t.Run("unsuccessful login due to invalid password", func(t *testing.T) {
		_, _, err = authService.Login(dto.LoginRequestDto{
			Login:    user.Nickname,
			Password: "invalid_password",
		}, gofakeit.UUID())
		assert.Error(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
}

func TestRefresh(t *testing.T) {
	var log = nop.GetLogger()
	var err error
	err = os.Setenv("ACCESS_TOKEN_EXP", "300")
	err = os.Setenv("REFRESH_TOKEN_EXP", "1209600")
	// mocks
	userRepo := new(mocks.UserRepository)
	tokenRepo := new(mocks.RefreshTokenRepository)
	resultsMigrator := new(mocks.ResultsMigrator)

	password := gofakeit.Password(true, true, true, true, false, 8)
	hashPassword, err := HashPassword(password)
	user := domain.User{
		Nickname: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: hashPassword,
	}

	refresh, err := utils.GenerateRefreshJWT(user.ID.Hex())
	tokenRepo.
		On("GetRefreshToken", refresh).
		Return(domain.RefreshToken{
			User:  user.ID,
			Token: refresh,
		}, nil)
	tokenRepo.
		On("GetRefreshToken", mock.AnythingOfType("string")).
		Return(domain.RefreshToken{}, &errors.RefreshError{})
	tokenRepo.
		On("UpdateRefreshToken", refresh, mock.Anything).
		Return(domain.RefreshToken{}, nil)

	userRepo.
		On("GetUserByID", user.ID.Hex()).
		Return(user, nil)

	// service
	authService := NewAuthService(userRepo, tokenRepo, resultsMigrator, log)

	t.Run("successful refresh", func(t *testing.T) {
		_, _, err = authService.Refresh(refresh)
		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
	t.Run("unsuccessful refresh due to invalid refresh token", func(t *testing.T) {
		_, _, err = authService.Refresh("invalid_refresh_token")
		assert.Error(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
}

func TestLogout(t *testing.T) {
	var log = nop.GetLogger()
	// mocks
	userRepo := new(mocks.UserRepository)
	tokenRepo := new(mocks.RefreshTokenRepository)
	resultsMigrator := new(mocks.ResultsMigrator)

	refresh, err := utils.GenerateRefreshJWT(gofakeit.UUID())
	tokenRepo.
		On("DeleteRefreshToken", refresh).
		Return(nil)
	tokenRepo.
		On("DeleteRefreshToken", mock.AnythingOfType("string")).
		Return(&errors.RefreshError{})

	// service
	authService := NewAuthService(userRepo, tokenRepo, resultsMigrator, log)

	t.Run("successful logout", func(t *testing.T) {
		authService.Logout(refresh)
		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
		resultsMigrator.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
}
