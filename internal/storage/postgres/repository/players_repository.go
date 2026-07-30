package postgres_repository

import (
	"context"
	"database/sql"

	"github.com/1tzArad/wyrm/internal/player"
	sqlc "github.com/1tzArad/wyrm/internal/storage/postgres/generated"
	"github.com/google/uuid"
)

type PlayerRepository struct {
	queries *sqlc.Queries
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		queries: sqlc.New(db),
	}
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
