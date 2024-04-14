package mongodb

import (
	"code-typing-auth-service/internal/core/domain"
	"code-typing-auth-service/internal/core/ports"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type RefreshTokenRepository struct {
}

func NewRefreshTokenRepository() ports.RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

func (r *RefreshTokenRepository) GetRefreshToken(token string) (refreshToken domain.RefreshToken, err error) {
	err = mgm.Coll(&refreshToken).First(bson.M{"token": token}, &refreshToken)
	if err != nil {
		return refreshToken, fmt.Errorf("refresh token '%s' not found", token)
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepository) CreateRefreshToken(refreshToken domain.RefreshToken) (ID string, err error) {
	err = mgm.Coll(&refreshToken).Create(&refreshToken)
	if err != nil {
		err = fmt.Errorf(`token not created due to error: %v`, err)
		return
	}
	return refreshToken.ID.Hex(), nil
}

func (r *RefreshTokenRepository) UpdateRefreshToken(oldRefreshToken, newRefreshToken string) (refreshToken domain.RefreshToken, err error) {
	refreshToken, _ = r.GetRefreshToken(oldRefreshToken)
	refreshToken.Token = newRefreshToken
	err = mgm.Coll(&refreshToken).Update(&refreshToken)
	if err != nil {
		return refreshToken, fmt.Errorf(`token not updated due to error: %v`, err)
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(token string) (err error) {
	var refreshToken domain.RefreshToken
	err = mgm.Coll(&refreshToken).First(bson.M{"token": token}, &refreshToken)
	if err != nil {
		return fmt.Errorf("refresh token '%s' not found", token)
	}
	err = mgm.Coll(&refreshToken).Delete(&refreshToken)
	if err != nil {
		return fmt.Errorf(`token not deleted due to error: %v`, err)
	}
	return nil
}
