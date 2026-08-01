package game

import (
	"context"

	"github.com/1tzArad/wyrm/internal/player"
	"github.com/google/uuid"
)

type PlayerStore interface {
	SavePlayer(ctx context.Context, p *player.Player) error
	LoadPlayer(ctx context.Context, id uuid.UUID) (*player.Player, error)
	GetByID(ctx context.Context, id uuid.UUID) (*player.Player, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*player.Player, error)
	Create(ctx context.Context, user_uud uuid.UUID) (*player.Player, error)
}
