package service

import (
	"context"
	"fmt"
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/model"
	"go-my-demo/internal/repository"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*model.User, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	GetAllUsers(req v1.GetAllUsersRequest, ctx context.Context) ([]model.User, error)
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	// userId, err := s.sid.GenString()
	// if err != nil {
	// 	return err
	// }
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return "", v1.ErrLoginFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	fmt.Println("err:", err)
	if err != nil {
		return "", err
	}

	// 将 user.Id 从 uint64 转换为 string
	userIdStr := strconv.FormatUint(user.Id, 10)
	// token有效期为1天
	token, err := s.jwt.GenToken(userIdStr, time.Now().Add(time.Hour*24*1))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Nickname = req.Nickname

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAllUsers(req v1.GetAllUsersRequest, ctx context.Context) ([]model.User, error) {
	users, err := s.userRepo.GetAll(req, ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
