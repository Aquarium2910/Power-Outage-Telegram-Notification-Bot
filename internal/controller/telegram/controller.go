package telegram

import (
	"NEW_LOE_bot/internal/usecase"
	"NEW_LOE_bot/pkg/logger"
)

type Controller struct {
	uc usecase.User
	l  logger.Interface
}
