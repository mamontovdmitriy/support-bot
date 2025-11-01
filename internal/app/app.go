package app

import (
	"fmt"
	"support-bot/config"

	log "github.com/sirupsen/logrus"
)

func Run(configPath string) {
	log.Info("Init application...")

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - config error: %w", err))
	}

	log.Infof("%v", cfg)
}
