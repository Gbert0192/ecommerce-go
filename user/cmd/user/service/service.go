package service

import (
	"context"
	"user/cmd/user/repository"
	"user/models"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (svc *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := svc.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (svc *UserService) CreateNewUser(ctx context.Context, user *models.User) (int64, error) {
	userID, err := svc.UserRepo.CreateNewUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (svc *UserService) GetUserById(ctx context.Context, userId int64) (*models.User, error) {
	user, err := svc.UserRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil

}
