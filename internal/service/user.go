package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/dwikie/auth-manager/internal/domain"
	"golang.org/x/crypto/argon2"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *domain.User, profile *domain.UserProfile) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserById(ctx context.Context, username string) (*domain.User, error)
	GetUserProfileByUsername(ctx context.Context, username string) (*domain.UserProfile, error)
	GetUserProfileByUserId(ctx context.Context, userid int) (*domain.UserProfile, error)
}

type UserService struct {
	userRepository IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
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

	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username:  req.Username,
		Password:  hashedPassword,
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

func (s *UserService) GetUserById(ctx context.Context, userid string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Printf("Getting user by id: %s", userid)

	return s.userRepository.GetUserById(ctx, userid)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Printf("Getting user by username: %s", username)

	return s.userRepository.GetUserByUsername(ctx, username)
}

func (s *UserService) GetUserProfileByUsername(ctx context.Context, username string) (*domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Printf("Getting user profile by username: %s", username)

	return s.userRepository.GetUserProfileByUsername(ctx, username)
}

func (s *UserService) GetUserProfileById(ctx context.Context, userid int) (*domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Printf("Getting user profile by id: %d", userid)

	return s.userRepository.GetUserProfileByUserId(ctx, userid)
}

func (s *UserService) generateSalt(n int) ([]byte, error) {
	salt := make([]byte, n)
	_, err := rand.Read(salt)
	return salt, err
}

func (s *UserService) HashPassword(password string) (encodedHash string, err error) {
	salt, err := s.generateSalt(16)
	if err != nil {
		return "", err
	}

	var time uint32 = 1
	var memory uint32 = 64 * 1024
	var threads uint8 = 4
	var keyLen uint32 = 32

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, time, threads, b64Salt, b64Hash)

	return encodedHash, nil
}

func (s *UserService) VerifyPassword(encodedHash, password string) bool {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return false
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	return subtle.ConstantTimeCompare(computedHash, hash) == 1
}
