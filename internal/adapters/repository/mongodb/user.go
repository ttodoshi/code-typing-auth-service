package mongodb

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/domain"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
}

func NewUserRepository() ports.UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUserByID(ID string) (user domain.User, err error) {
	err = mgm.Coll(&user).FindByID(ID, &user)
	if err != nil {
		return user, fmt.Errorf("user by ID '%s' not found", ID)
	}
	return user, nil
}

func (r *UserRepository) GetUserByNickname(nickname string) (user domain.User, err error) {
	err = mgm.Coll(&user).First(bson.M{"nickname": nickname}, &user)
	if err != nil {
		return user, fmt.Errorf("user by nickname '%s' not found", nickname)
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (user domain.User, err error) {
	err = mgm.Coll(&user).First(bson.M{"email": email}, &user)
	if err != nil {
		return user, fmt.Errorf("user by email '%s' not found", email)
	}
	return user, nil
}

func (r *UserRepository) SaveUser(user domain.User) (domain.User, error) {
	err := mgm.Coll(&user).Create(&user)
	if err != nil {
		return user, fmt.Errorf(`user not created due to error: %v`, err)
	}
	return user, nil
}
