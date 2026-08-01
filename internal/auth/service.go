package auth

import (
	"context"
	"errors"

	sqlc "github.com/1tzArad/wyrm/internal/storage/postgres/generated"
)

// Errors

var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrUserExists = errors.New("username already exists")

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (sqlc.User, bool, error)
	Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
}

type Service struct {
	userRepo UserRepository
}

func NewService(userRepo UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Register(ctx context.Context, username, password string) error {
	_, found, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if found {
		return ErrUserExists
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.userRepo.Create(ctx, sqlc.CreateUserParams{Username: username, PasswordHash: hashedPassword})
	return err
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, found, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil || !found {
		return "", ErrInvalidCredentials
	}

	if !checkPassword(user.PasswordHash, password) {
		return "", ErrInvalidCredentials
	}

	return GenerateJWT(user.ID)
}
