package game

import (
	"context"

	"github.com/1tzArad/wyrm/internal/player"
	sqlc "github.com/1tzArad/wyrm/internal/storage/postgres/generated"
	"github.com/google/uuid"
)

type PlayerStore interface {
	SavePlayer(ctx context.Context, p *player.Player) error
	LoadPlayer(ctx context.Context, id uuid.UUID) (*player.Player, error)
	GetByID(ctx context.Context, id uuid.UUID) (sqlc.Player, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (sqlc.Player, error)
	Create(ctx context.Context, arg sqlc.CreatePlayerParams) (sqlc.Player, error)
}
