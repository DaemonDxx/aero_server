package main

import (
	"fmt"
	"github.com/daemondxx/lks_back/internal/app"
	"github.com/daemondxx/lks_back/internal/config"
	"github.com/daemondxx/lks_back/internal/logger"
)

// @title			AeroBot server
// @version		1.0
// @description	Server for aero bot
func main() {

	log := logger.NewLogger(logger.DEV)

	log.Info().Msg("init config...")
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("init config error: %e", err))
	}
	log.Info().Msg("config init successful")

	log.Info().Msg("init aeroserver...")
	a, err := app.NewApp(cfg, log)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("init aeroserver error: %e", err))
	}
	log.Info().Msg("aeroserver init successful")

	log.Info().Msg(fmt.Sprintf("start aeroserver on port %s", cfg.Http.Port))
	if err := a.Run(); err != nil {
		log.Fatal().Msg(fmt.Sprintf("start aeroserver error: %e", err))
	}

}
