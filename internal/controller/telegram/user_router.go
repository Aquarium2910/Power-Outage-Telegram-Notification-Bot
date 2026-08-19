package telegram

import (
	"NEW_LOE_bot/internal/usecase"
	"NEW_LOE_bot/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

func NewRoutes(g *tele.Group, uc usecase.User, l logger.Interface) {
	h := &Controller{
		uc: uc,
		l:  l,
	}

	g.Handle("/start", h.StartHandler)
	g.Handle(tele.OnMyChatMember, h.OnMyChatMember)
}
