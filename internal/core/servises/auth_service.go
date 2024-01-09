package servises

import (
	"github.com/jinzhu/copier"
	"speed-typing-auth-service/internal/adapters/dto"
	"speed-typing-auth-service/internal/core/domain"
	"speed-typing-auth-service/internal/core/errors"
	"speed-typing-auth-service/internal/core/ports"
	"speed-typing-auth-service/internal/core/utils"
	"speed-typing-auth-service/pkg/logging"
)

type AuthService struct {
	userRepo  ports.UserRepository
	tokenRepo ports.RefreshTokenRepository
	log       logging.Logger
}

func NewAuthService(userRepo ports.UserRepository, tokenRepo ports.RefreshTokenRepository, log logging.Logger) ports.AuthService {
	return &AuthService{userRepo: userRepo, tokenRepo: tokenRepo, log: log}
}

func (s *AuthService) Register(registerRequestDto dto.RegisterRequestDto) (accessToken string, refreshToken string, err error) {
	var user domain.User

	registerRequestDto.Password, err = utils.HashPassword(registerRequestDto.Password)
	if err != nil {
		return accessToken, refreshToken, err
	}

	err = copier.Copy(&user, &registerRequestDto)
	user, err = s.userRepo.SaveUser(user)
	if err != nil {
		return accessToken, refreshToken, err
	}

	accessToken, err = utils.GenerateAccessJWT(user)
	refreshToken, err = utils.GenerateRefreshJWT(user.ID.Hex())
	if err != nil {
		return accessToken, refreshToken, err
	}

	_, err = s.tokenRepo.SaveRefreshToken(domain.RefreshToken{
		User:  user.ID,
		Token: refreshToken,
	})
	if err != nil {
		return accessToken, refreshToken, err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) Login(loginRequestDto dto.LoginRequestDto) (accessToken string, refreshToken string, err error) {
	var user domain.User
	user, err = s.userRepo.GetUserByNickname(loginRequestDto.Login)
	if err != nil {
		user, err = s.userRepo.GetUserByEmail(loginRequestDto.Login)
		if err != nil {
			return accessToken, refreshToken, err
		}
	}

	err = utils.VerifyPassword(user.Password, loginRequestDto.Password)
	if err != nil {
		return accessToken, refreshToken, &errors.LoginOrPasswordDoNotMatchError{
			Message: "login or password do not match",
		}
	}

	accessToken, err = utils.GenerateAccessJWT(user)
	refreshToken, err = utils.GenerateRefreshJWT(user.ID.Hex())
	if err != nil {
		return accessToken, refreshToken, err
	}

	_, err = s.tokenRepo.SaveRefreshToken(domain.RefreshToken{
		User:  user.ID,
		Token: refreshToken,
	})
	if err != nil {
		return accessToken, refreshToken, err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(refreshRequestDto dto.RefreshRequestDto) (accessToken string, refreshToken string, err error) {
	token, err := s.tokenRepo.GetRefreshToken(refreshRequestDto.RefreshToken)
	if err != nil {
		return accessToken, refreshToken, err
	}

	user, _ := s.userRepo.GetUserByID(token.User.Hex())

	accessToken, err = utils.GenerateAccessJWT(user)
	refreshToken, err = utils.GenerateRefreshJWT(user.ID.Hex())
	if err != nil {
		return accessToken, refreshToken, err
	}

	_, err = s.tokenRepo.UpdateRefreshToken(token.Token, refreshToken)
	if err != nil {
		return accessToken, refreshToken, err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) Logout(logoutRequestDto dto.LogoutRequestDto) {
	_ = s.tokenRepo.DeleteRefreshToken(logoutRequestDto.RefreshToken)
}
