package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

var (
	secretKey          = []byte(os.Getenv("SECRET_KEY"))
	AccessTokenExp, _  = strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	RefreshTokenExp, _ = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))
)

type Claim struct {
	Name  string
	Value interface{}
}

func GenerateAccessJWT(sub string, claims ...Claim) (accessToken string, err error) {
	accessToken, err = generateJWT(sub, AccessTokenExp, claims...)

	if err != nil {
		err = fmt.Errorf("access jwt generation error due to: %s", err.Error())
		return
	}
	return
}

func GenerateRefreshJWT(sub string, claims ...Claim) (refreshToken string, err error) {
	refreshToken, err = generateJWT(sub, RefreshTokenExp, claims...)

	if err != nil {
		err = fmt.Errorf("refresh jwt generation error due to: %s", err.Error())
		return
	}
	return
}

func generateJWT(sub string, exp int, claims ...Claim) (jwtToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	tokenClaims["sub"] = sub
	for _, claim := range claims {
		tokenClaims[claim.Name] = claim.Value
	}
	tokenClaims["iat"] = time.Now().Unix()
	tokenClaims["exp"] = time.Now().Unix() + int64(exp)

	jwtToken, err = token.SignedString(secretKey)
	return
}
