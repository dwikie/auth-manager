package service

import (
	"context"
	"time"

	"github.com/dwikie/auth-manager/internal/domain"
)

type UserRepository interface {
	CreateUser(context.Context, *domain.User, *domain.UserProfile) error
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

type CreateUserRequest struct {
	Username  string
	Password  string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &domain.User{
		Username:  req.Username,
		Password:  req.Password,
		Status:    0, // Set the initial status to 0 (inactive)
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	profile := &domain.UserProfile{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		EmailStatus: 0, // Set the initial email status to 0 (unverified)
		Phone:       req.Phone,
		PhoneStatus: 0, // Set the initial phone status to 0 (unverified)
	}

	return s.userRepository.CreateUser(ctx, user, profile)
}
