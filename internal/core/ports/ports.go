package ports

import (
	"github.com/ttodoshi/code-typing-auth-service/internal/core/domain"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/ports/dto"
)

type AuthService interface {
	Register(registerRequestDto dto.RegisterRequestDto, session string) (access string, refresh string, err error)
	Login(loginRequestDto dto.LoginRequestDto, session string) (access string, refresh string, err error)
	Refresh(oldRefreshToken string) (access string, refresh string, err error)
	Logout(refreshToken string)
}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=RefreshTokenRepository
type RefreshTokenRepository interface {
	GetRefreshToken(refreshToken string) (domain.RefreshToken, error)
	CreateRefreshToken(refreshToken domain.RefreshToken) (string, error)
	UpdateRefreshToken(oldRefreshToken, newRefreshToken string) (domain.RefreshToken, error)
	DeleteRefreshToken(refreshToken string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=UserRepository
type UserRepository interface {
	GetUserByID(ID string) (domain.User, error)
	GetUserByNickname(nickname string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	SaveUser(user domain.User) (domain.User, error)
}

const (
	AuthExchange = "auth-exchange"
)

var Exchanges = []string{AuthExchange}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=EventDispatcher
type EventDispatcher interface {
	Dispatch(event domain.Event)
}
