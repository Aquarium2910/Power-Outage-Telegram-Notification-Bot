package usecase

import (
	"NEW_LOE_bot/internal/entity"
	"context"
)

type (
	User interface {
		Register(ctx context.Context, telegramID int64, username string, firstName string) (entity.User, error)
		GetByTelegramID(ctx context.Context, telegramID int64) (entity.User, error)
		UpdateIsActive(ctx context.Context, telegramID int64, status bool) error
	}
)
