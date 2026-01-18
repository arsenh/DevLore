package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arsenh/DevLore/internal/model"
	"github.com/arsenh/DevLore/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository model.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: repository.NewDummyUserRepository(),
	}
}

func (u *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	//TODO: handle FindByID error on real database to recognize DB error
	user, _ := u.userRepository.FindByID(ctx, id)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to find user by id: %d", id)
	//}
	return user, nil
}

func (u *UserService) RegisterUser(ctx context.Context, email string, password string, fullName string) (*model.User, error) {

	user := &model.User{
		ID:        uuid.New(),
		FullName:  fullName,
		Email:     email,
		Password:  password, // TODO: need to create HASH on password
		CreatedAt: time.Now(),
	}

	err := u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to register user by email: %s", email)
	}

	return user, nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) *model.User {
	user := u.userRepository.FindByEmail(ctx, email)
	return user
}

func (u *UserService) CheckUserPassword(ctx context.Context, userID uuid.UUID, inputPassword string) bool {
	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		return false
	}

	//TODO: In real database, need to compare hashes, from user input and database hash
	return user.Password == inputPassword
}
