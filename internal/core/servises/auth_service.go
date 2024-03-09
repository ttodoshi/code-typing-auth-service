package servises

import (
	"code-typing-auth-service/internal/adapters/dto"
	"code-typing-auth-service/internal/core/domain"
	"code-typing-auth-service/internal/core/errors"
	"code-typing-auth-service/internal/core/ports"
	"code-typing-auth-service/internal/core/utils"
	"code-typing-auth-service/pkg/logging"
	"github.com/jinzhu/copier"
)

type AuthService struct {
	userRepo        ports.UserRepository
	tokenRepo       ports.RefreshTokenRepository
	resultsMigrator ports.ResultsMigrator
	log             logging.Logger
}

func NewAuthService(userRepo ports.UserRepository, tokenRepo ports.RefreshTokenRepository, resultsMigrator ports.ResultsMigrator, log logging.Logger) ports.AuthService {
	return &AuthService{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		resultsMigrator: resultsMigrator,
		log:             log,
	}
}

func (s *AuthService) Register(registerRequestDto dto.RegisterRequestDto, session string) (access string, refresh string, err error) {
	var user domain.User

	registerRequestDto.Password, err = utils.HashPassword(registerRequestDto.Password)
	if err != nil {
		return
	}

	err = copier.Copy(&user, &registerRequestDto)
	user, err = s.userRepo.SaveUser(user)
	if err != nil {
		return
	}

	access, refresh, err = s.generateTokens(user)
	if err != nil {
		return
	}

	s.resultsMigrator.MigrateSessionResults(session, user.ID.Hex())
	_, err = s.tokenRepo.SaveRefreshToken(domain.RefreshToken{
		User:  user.ID,
		Token: refresh,
	})
	return
}

func (s *AuthService) Login(loginRequestDto dto.LoginRequestDto, session string) (access string, refresh string, err error) {
	var user domain.User
	user, err = s.userRepo.GetUserByNickname(loginRequestDto.Login)
	if err != nil {
		user, err = s.userRepo.GetUserByEmail(loginRequestDto.Login)
		if err != nil {
			return
		}
	}

	err = utils.VerifyPassword(user.Password, loginRequestDto.Password)
	if err != nil {
		return access, refresh, &errors.LoginOrPasswordDoNotMatchError{
			Message: "login or password do not match",
		}
	}

	access, refresh, err = s.generateTokens(user)
	if err != nil {
		return
	}

	s.resultsMigrator.MigrateSessionResults(session, user.ID.Hex())
	_, err = s.tokenRepo.SaveRefreshToken(domain.RefreshToken{
		User:  user.ID,
		Token: refresh,
	})
	return
}

func (s *AuthService) Refresh(oldRefreshToken string) (access string, refresh string, err error) {
	token, err := s.tokenRepo.GetRefreshToken(oldRefreshToken)
	if err != nil {
		return
	}

	user, _ := s.userRepo.GetUserByID(token.User.Hex())

	access, refresh, err = s.generateTokens(user)
	if err != nil {
		return
	}

	_, err = s.tokenRepo.UpdateRefreshToken(token.Token, refresh)
	return
}

func (s *AuthService) generateTokens(user domain.User) (accessToken string, refreshToken string, err error) {
	accessToken, err = utils.GenerateAccessJWT(user)
	refreshToken, err = utils.GenerateRefreshJWT(user.ID.Hex())
	return
}

func (s *AuthService) Logout(refreshToken string) {
	_ = s.tokenRepo.DeleteRefreshToken(refreshToken)
}
