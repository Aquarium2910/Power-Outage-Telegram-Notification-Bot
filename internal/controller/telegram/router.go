package telegram

import (
	"NEW_LOE_bot/internal/usecase"
	"NEW_LOE_bot/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

func NewRouter(g *tele.Group, uc usecase.User, l logger.Interface) {

	NewRoutes(g, uc, l)
}
