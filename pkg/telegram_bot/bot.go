package telegram_bot

import (
	"NEW_LOE_bot/pkg/logger"
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	ctx    context.Context
	eg     *errgroup.Group
	tb     *telebot.Bot
	l      logger.Interface
	notify chan error
}

func New(token string, poller telebot.Poller, l logger.Interface) (*Bot, error) {
	settings := telebot.Settings{
		Token:   token,
		Poller:  poller,
		Offline: true,
	}

	// TODO: Add options for settings

	bot, err := telebot.NewBot(settings)
	if err != nil {
		l.Fatal("Failed to create telegram_bot bot", err)
		return nil, err
	}

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1)

	return &Bot{
		ctx:    ctx,
		eg:     group,
		tb:     bot,
		l:      l,
		notify: make(chan error, 1),
	}, nil
}

func (b *Bot) Start() {
	b.eg.Go(func() error {
		b.tb.Start()
		return nil
	})

	b.l.Info("telegram_bot - Start - Telegram bot started")
}

func (b *Bot) Shutdown() error {
	b.tb.Stop()

	err := b.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		b.l.Error(err, "telegram_bot bot - bot - Shutdown - b.eg.Wait()")

		return err
	}

	return nil
}

func (b *Bot) Notify() <-chan error {
	return b.notify
}

func (b *Bot) GetBot() *telebot.Bot {
	return b.tb
}
