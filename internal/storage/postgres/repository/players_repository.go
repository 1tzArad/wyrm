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

func (r *PlayerRepository) Create(ctx context.Context, user_uuid uuid.UUID) (*player.Player, error) {
	row, err := r.queries.CreatePlayer(ctx, user_uuid)
	if err != nil {
		return nil, err
	}
	return toPlayer(row), nil
}

func (r *PlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (*player.Player, error) {
	row, err := r.queries.GetPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	return toPlayer(row), nil
}

func (r *PlayerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*player.Player, error) {
	row, err := r.queries.GetPlayerByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	return toPlayer(row), nil
}

func (r *PlayerRepository) GetPlayersByUserID(ctx context.Context, userID uuid.UUID) ([]*player.Player, error) {
	rows, err := r.queries.GetPlayersByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	players := make([]*player.Player, 0, len(rows))
	for _, row := range rows {
		players = append(players, toPlayer(row))
	}
	return players, nil
}

func (r *PlayerRepository) GetAll(ctx context.Context) ([]*player.Player, error) {
	rows, err := r.queries.GetAllPlayers(ctx)
	if err != nil {
		return nil, err
	}
	players := make([]*player.Player, 0, len(rows))
	for _, row := range rows {
		players = append(players, toPlayer(row))
	}
	return players, nil
}

func toPlayer(row sqlc.Player) *player.Player {
	return &player.Player{
		ID:     row.ID,
		X:      row.X,
		Y:      row.Y,
		Health: row.Hp,
		Mana:   row.Mana,
	}
}

func (r *PlayerRepository) LoadPlayer(ctx context.Context, id uuid.UUID) (*player.Player, error) {
	row, err := r.queries.GetPlayer(ctx, id)
	if err != nil {
		return nil, err
	}
	return toPlayer(row), nil
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
