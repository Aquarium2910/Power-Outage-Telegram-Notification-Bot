package main

import (
	"NEW_LOE_bot/config"
	"NEW_LOE_bot/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.RunMigrations(cfg)
	app.Run(cfg)
}
