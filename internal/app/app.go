package app

import (
	"NEW_LOE_bot/config"
	"NEW_LOE_bot/internal/controller/telegram"
	"NEW_LOE_bot/internal/repo/persistent"
	"NEW_LOE_bot/internal/usecase/user"
	"NEW_LOE_bot/pkg/logger"
	"NEW_LOE_bot/pkg/postgres"
	"NEW_LOE_bot/pkg/telegram_bot"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/telebot.v3"
)

type useCases struct {
	user *user.UseCase
}

type servers struct {
	telegram *telegram_bot.Bot
}

func initUseCases(pg *postgres.Postgres) useCases {
	userRepo := persistent.NewUserRepo(pg)

	return useCases{
		user: user.New(userRepo),
	}
}

func initServers(cfg *config.Config, uc useCases, l logger.Interface) servers {
	poller := &telebot.LongPoller{
		Limit:   cfg.Telegram.PollerLimit,
		Timeout: cfg.Telegram.PollerTimeout,
	}

	tg, err := telegram_bot.New(cfg.Telegram.Token, poller, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - telegram_bot - telegram_bot.New: %w", err))
	}

	telegram.NewRouter(tg.GetBot().Group(), uc.user, l)

	return servers{
		telegram: tg,
	}
}

func (s *servers) startServers() {
	s.telegram.Start()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info(fmt.Sprintf("app - Run - signal: %s", sig.String()))
	case err = <-s.telegram.Notify():
		l.Error(fmt.Errorf("app - Run - telegram.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.telegram.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - telegram.Shutdown: %w", err))
	}
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL(), postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	uc := initUseCases(pg)
	s := initServers(cfg, uc, l)
	s.startServers()
	s.waitForShutdown(l)
}
