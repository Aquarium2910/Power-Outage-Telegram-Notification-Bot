package repo

import (
	"NEW_LOE_bot/internal/entity"
	"context"
)

type (
	UserRepo interface {
		Store(ctx context.Context, user *entity.User) error
		GetByTelegramID(ctx context.Context, telegramID int64) (entity.User, error)
		Update(ctx context.Context, user *entity.User) error
	}
)
