package telegram

import (
	"NEW_LOE_bot/internal/entity"
	"context"
	"errors"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (r *Controller) StartHandler(c tele.Context) error {
	tgUser := c.Sender()
	if tgUser == nil {
		return nil
	}

	ctx := context.Background()

	_, err := r.uc.Register(ctx, tgUser.ID, tgUser.Username, tgUser.FirstName)

	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			r.l.Error(err, "telegram_bot - StartHandler")

			//TODO: Add presenter
			return c.Send("Ви уже зареєстровані!")
		}

		r.l.Error(err, "telegram_bot - StartHandler")

		return c.Send("Сталась невідома помилка при реєстрації. Спробуйте ще раз пізніше")
	}

	//TODO: Add presenter
	return c.Send(fmt.Sprintf("Вітаю, %s. *Решта повідомлення як побачу в інтерфейсі що треба підказати*", tgUser.FirstName))

}

func (r *Controller) OnMyChatMember(c tele.Context) error {
	update := c.ChatMember()
	if update == nil {
		return nil
	}

	if update.Chat.Type != tele.ChatPrivate {
		return nil
	}

	telegramID := update.Chat.ID
	newStatus := update.NewChatMember.Role

	ctx := context.Background()

	switch newStatus {
	case tele.Kicked:
		r.l.Info(fmt.Sprintf("User %d blocked bot", telegramID))
		err := r.uc.UpdateIsActive(ctx, telegramID, false)
		if err != nil {
			r.l.Error(err, "telegram_bot - OnMyChatMemberHandler - Kicked")
			return err
		}
	case tele.Member:
		r.l.Info(fmt.Sprintf("User %d unblocked bot", telegramID))
		err := r.uc.UpdateIsActive(ctx, telegramID, true)
		if err != nil {
			r.l.Error(err, "telegram_bot - OnMyChatMemberHandler - Member")
			return err
		}
	}

	return nil
}
