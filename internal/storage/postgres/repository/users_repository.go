package postgres_repository

import (
	"context"
	"database/sql"
	"errors"

	sqlc "github.com/1tzArad/wyrm/internal/storage/postgres/generated"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		queries: sqlc.New(db),
	}
}

func (r *UserRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
	user, err := r.queries.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, ErrUserNotFound
		}
		return sqlc.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (sqlc.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, ErrUserNotFound
		}
		return sqlc.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]sqlc.User, error) {
	return r.queries.GetAllUsers(ctx)
}

func (r *UserRepository) Update(ctx context.Context, arg sqlc.UpdateUserParams) error {
	return r.queries.UpdateUser(ctx, arg)
}
