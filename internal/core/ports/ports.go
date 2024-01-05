package ports

import (
	"speed-typing-auth-service/internal/adapters/dto"
	"speed-typing-auth-service/internal/core/domain"
)

type AuthService interface {
	Register(registerRequestDto dto.RegisterRequestDto) (accessToken string, refreshToken string, err error)
	Login(loginRequestDto dto.LoginRequestDto) (accessToken string, refreshToken string, err error)
	Refresh(refreshRequestDto dto.RefreshRequestDto) (accessToken string, refreshToken string, err error)
	Logout(logoutRequestDto dto.LogoutRequestDto)
}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=RefreshTokenRepository
type RefreshTokenRepository interface {
	GetRefreshToken(refreshToken string) (domain.RefreshToken, error)
	SaveRefreshToken(refreshToken domain.RefreshToken) (string, error)
	UpdateRefreshToken(oldRefreshToken, newRefreshToken string) (domain.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=UserRepository
type UserRepository interface {
	GetUserByID(ID string) (domain.User, error)
	GetUserByNickname(nickname string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	SaveUser(user domain.User) (string, error)
}
