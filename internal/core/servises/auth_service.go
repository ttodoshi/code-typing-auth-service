package servises

import (
	"code-typing-auth-service/internal/core/domain"
	"code-typing-auth-service/internal/core/ports"
	"code-typing-auth-service/internal/core/ports/dto"
	"code-typing-auth-service/pkg/jwt"
	"code-typing-auth-service/pkg/logging"
	"code-typing-auth-service/pkg/password"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
)

type AuthService struct {
	userRepo        ports.UserRepository
	tokenRepo       ports.RefreshTokenRepository
	eventDispatcher ports.EventDispatcher
	log             logging.Logger
}

func NewAuthService(userRepo ports.UserRepository, tokenRepo ports.RefreshTokenRepository, resultsMigrator ports.EventDispatcher, log logging.Logger) ports.AuthService {
	return &AuthService{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		eventDispatcher: resultsMigrator,
		log:             log,
	}
}

func (s *AuthService) Register(registerRequestDto dto.RegisterRequestDto, session string) (access string, refresh string, err error) {
	var user domain.User

	registerRequestDto.Password, err = password.HashPassword(registerRequestDto.Password)
	if err != nil {
		return
	}

	err = copier.Copy(&user, &registerRequestDto)
	if err != nil {
		err = fmt.Errorf(`struct mapping error: %w`, ports.InternalServerError)
		return
	}

	user, err = s.saveUser(user)
	if err != nil {
		return
	}

	access, refresh, err = s.generateTokens(user)
	if err != nil {
		return
	}

	_, err = s.tokenRepo.CreateRefreshToken(domain.RefreshToken{
		User:  user.ID,
		Token: refresh,
	})
	if err != nil {
		err = fmt.Errorf(`creating refresh token error: %w`, ports.InternalServerError)
		return
	}
	s.dispatchAuthEvent(session, user)
	return
}

func (s *AuthService) saveUser(user domain.User) (domain.User, error) {
	var err error
	_, err = s.userRepo.GetUserByNickname(user.Nickname)
	if err == nil {
		err = fmt.Errorf("nickname already picked: %w", ports.BadRequestError)
		return domain.User{}, err
	}
	_, err = s.userRepo.GetUserByEmail(user.Email)
	if err == nil {
		err = fmt.Errorf("account with this email already exists: %w", ports.BadRequestError)
		return domain.User{}, err
	}

	user, err = s.userRepo.SaveUser(user)
	if err != nil {
		s.log.Warnf("user not saved due to error: %v", err)
		err = fmt.Errorf(`saving user error: %w`, ports.InternalServerError)
		return domain.User{}, err
	}
	return user, nil
}

func (s *AuthService) Login(loginRequestDto dto.LoginRequestDto, session string) (access string, refresh string, err error) {
	var user domain.User
	user, err = s.userRepo.GetUserByNickname(loginRequestDto.Login)
	if err != nil {
		user, err = s.userRepo.GetUserByEmail(loginRequestDto.Login)
		if err != nil {
			err = fmt.Errorf("user not found: %w", ports.BadRequestError)
			return
		}
	}

	err = password.VerifyPassword(user.Password, loginRequestDto.Password)
	if err != nil {
		return access, refresh, fmt.Errorf(
			"login or password do not match: %w", ports.BadRequestError,
		)
	}

	access, refresh, err = s.generateTokens(user)
	if err != nil {
		return
	}

	_, err = s.tokenRepo.CreateRefreshToken(domain.RefreshToken{
		User:  user.ID,
		Token: refresh,
	})
	if err != nil {
		err = fmt.Errorf(`creating refresh token error: %w`, ports.InternalServerError)
	}
	s.dispatchAuthEvent(session, user)
	return
}

func (s *AuthService) dispatchAuthEvent(session string, user domain.User) {
	if session == "" {
		return
	}

	body, err := json.Marshal(
		map[string]interface{}{
			"session": session,
			"userID":  user.ID.Hex(),
		},
	)
	if err != nil {
		s.log.Warnf(`error marshaling event body: %v`, err)
		return
	}

	s.eventDispatcher.Dispatch(domain.Event{
		Exchange: ports.AuthExchange,
		Body:     body,
	})
}

func (s *AuthService) Refresh(oldRefreshToken string) (access string, refresh string, err error) {
	token, err := s.tokenRepo.GetRefreshToken(oldRefreshToken)
	if err != nil {
		err = fmt.Errorf("refresh token not found: %w", ports.UnauthorizedError)
		return
	}

	user, _ := s.userRepo.GetUserByID(token.User.Hex())

	access, refresh, err = s.generateTokens(user)
	if err != nil {
		return
	}

	_, err = s.tokenRepo.UpdateRefreshToken(token.Token, refresh)
	if err != nil {
		err = fmt.Errorf(`updating refresh token error: %w`, ports.InternalServerError)
	}
	return
}

func (s *AuthService) generateTokens(user domain.User) (accessToken string, refreshToken string, err error) {
	accessToken, err = jwt.GenerateAccessJWT(
		user.ID.Hex(),
		jwt.Claim{
			Name:  "nickname",
			Value: user.Nickname,
		},
	)
	if err != nil {
		err = fmt.Errorf(`generating tokens error: %w`, ports.InternalServerError)
		return
	}
	refreshToken, err = jwt.GenerateRefreshJWT(user.ID.Hex())
	if err != nil {
		err = fmt.Errorf(`generating tokens error: %w`, ports.InternalServerError)
		return
	}
	return
}

func (s *AuthService) Logout(refreshToken string) {
	err := s.tokenRepo.DeleteRefreshToken(refreshToken)
	if err != nil {
		s.log.Warnf("refresh token delete error: %v", err)
	}
}
