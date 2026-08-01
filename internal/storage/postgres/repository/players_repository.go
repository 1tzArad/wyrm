package postgres_repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/1tzArad/wyrm/internal/player"
	sqlc "github.com/1tzArad/wyrm/internal/storage/postgres/generated"
	"github.com/google/uuid"
)

var (
	ErrPlayerNotFound = errors.New("player not found")
)

type PlayerRepository struct {
	queries *sqlc.Queries
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		queries: sqlc.New(db),
	}
}

func (r *PlayerRepository) Create(ctx context.Context, arg sqlc.CreatePlayerParams) (sqlc.Player, error) {
	return r.queries.CreatePlayer(ctx, arg)
}

func (r *PlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (sqlc.Player, error) {
	row, err := r.queries.GetPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Player{}, ErrPlayerNotFound
		}
		return sqlc.Player{}, err
	}
	return row, nil
}

func (r *PlayerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (sqlc.Player, error) {
	row, err := r.queries.GetPlayerByUserId(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Player{}, ErrPlayerNotFound
		}
		return sqlc.Player{}, err
	}
	return row, nil
}

func (r *PlayerRepository) GetPlayersByUserID(ctx context.Context, userID uuid.UUID) ([]sqlc.Player, error) {
	return r.queries.GetPlayersByUserId(ctx, uuid.NullUUID{UUID: userID, Valid: true})
}

func (r *PlayerRepository) GetAll(ctx context.Context) ([]sqlc.Player, error) {
	return r.queries.GetAllPlayers(ctx)
}

func (r *PlayerRepository) LoadPlayer(ctx context.Context, id uuid.UUID) (*player.Player, error) {
	row, err := r.queries.GetPlayer(ctx, id)
	if err != nil {
		return nil, err
	}

	return &player.Player{
		ID:     row.ID,
		X:      row.X,
		Y:      row.Y,
		Health: row.Hp,
		Mana:   row.Mana,
	}, nil
}

func (r *PlayerRepository) SavePlayer(ctx context.Context, p *player.Player) error {
	return r.queries.SavePlayer(ctx, sqlc.SavePlayerParams{
		ID:   p.ID,
		X:    p.X,
		Y:    p.Y,
		Hp:   p.Health,
		Mana: p.Mana,
	})
}
