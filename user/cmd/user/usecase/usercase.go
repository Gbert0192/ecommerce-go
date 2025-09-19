package usecase

import (
	"context"
	"errors"
	"time"
	"user/cmd/user/service"
	"user/infrastructure/log"
	"user/models"
	"user/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUseCase struct {
	UserService service.UserService
	JWTSecret   string
}

func NewUserUseCase(userService service.UserService, secret string) *UserUseCase {
	return &UserUseCase{
		UserService: userService,
		JWTSecret:   secret,
	}
}

func (uc *UserUseCase) GetUserById(ctx context.Context, userId int64) (*models.User, error) {
	user, err := uc.UserService.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.User{}, nil
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := uc.UserService.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.User{}, nil
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": user.Email,
		}).Errorf("failed to hash password: %v", err)
		return err
	}
	user.Password = hashedPassword
	_, err = uc.UserService.CreateNewUser(ctx, user)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": user.Email,
			"name":  user.Name,
		}).Errorf("failed to create user: %v", err)
		return err
	}
	return nil
}

func (uc *UserUseCase) Login(ctx context.Context, param *models.LoginParameter) (string, error) {
	user, err := uc.UserService.UserRepo.GetUserByEmail(ctx, param.Email)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": param.Email,
		}).Errorf("GetUserByEmail got error: %v", err)
		return "", err
	}
	if user.ID == 0 {
		return "", errors.New("email not found")
	}
	isPasswordMatch := utils.CheckPasswordHash(param.Password, user.Password)

	if isPasswordMatch {
		return "", errors.New("wrong password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 6).Unix(),
	})
	tokenString, err := token.SignedString([]byte(uc.JWTSecret))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": param.Email,
		}).Errorf("GenerateToken Got Error: %v", err)
	}
	return tokenString, nil
}
