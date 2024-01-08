package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"speed-typing-auth-service/internal/core/domain"
	"speed-typing-auth-service/internal/core/errors"
	"strconv"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func GenerateAccessJWT(user domain.User) (accessToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	accessTokenExp, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXP"), 10, 64)
	if err != nil {
		return accessToken, &errors.TokenGenerationError{
			Message: fmt.Sprintf("access token generation error due to: %s", err.Error()),
		}
	}

	claims["ID"] = user.ID
	claims["nickname"] = user.Nickname
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Unix() + accessTokenExp

	accessToken, err = token.SignedString(secretKey)

	if err != nil {
		return accessToken, &errors.TokenGenerationError{
			Message: fmt.Sprintf("access token generation error due to: %s", err.Error()),
		}
	}
	return accessToken, nil
}

func GenerateRefreshJWT(ID string) (refreshToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	refreshTokenExp, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXP"), 10, 64)
	if err != nil {
		return refreshToken, &errors.TokenGenerationError{
			Message: fmt.Sprintf("refresh token generation error due to: %s", err.Error()),
		}
	}

	claims["ID"] = ID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Unix() + refreshTokenExp

	refreshToken, err = token.SignedString(secretKey)

	if err != nil {
		return refreshToken, &errors.TokenGenerationError{
			Message: fmt.Sprintf("refresh token generation error due to: %s", err.Error()),
		}
	}
	return refreshToken, nil
}

func IsTokenValid(token string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", &errors.TokenParsingError{
				Message: fmt.Sprintf("token parsing error"),
			}
		}
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, &errors.TokenParsingError{
			Message: "error while extracting claims from token",
		}
	}

	exp := claims["exp"].(int64)
	if exp < time.Now().Local().Unix() {
		return false, nil
	}
	return true, nil
}

func ExtractNickname(token string) (nickname string, err error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", &errors.TokenParsingError{
				Message: fmt.Sprintf("token parsing error"),
			}
		}
		return secretKey, nil
	})
	if err != nil {
		return nickname, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nickname, &errors.TokenParsingError{
			Message: "error while extracting claims from token",
		}
	}

	return claims["nickname"].(string), nil
}
