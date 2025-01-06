package app

import (
	"github.com/daemondxx/lks_back/internal/config"
	"github.com/daemondxx/lks_back/internal/dao"
	"github.com/daemondxx/lks_back/internal/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

type LKSApp struct {
	db   *gorm.DB
	log  *zerolog.Logger
	http *server.Server
}

func NewApp(cfg config.Config, log *zerolog.Logger) (*LKSApp, error) {
	app := &LKSApp{
		log: log,
	}

	log.Info().Msg("init database...")
	d, err := dao.NewDatabase(cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("init database failed")
		return nil, errors.Wrap(err, "failed to init database")
	}
	log.Info().Msg("database init successful")

	log.Info().Msg("start migrate database...")
	if err := d.AutoMigrate(); err != nil {
		log.Error().Err(err).Msg("failed to migrate database")
		return nil, errors.Wrap(err, "failed to migrate database")
	}
	log.Info().Msg("migrate database successfully")

	app.db = d.GetInstance()

	{
		l := log.With().Str("context", "http_server").Logger()
		app.http = server.NewServer(app.db, &cfg.Http, &l)
	}

	//if err := app.initServices(cfg); err != nil {
	//	return nil, fmt.Errorf("init services error: %e", err)
	//}

	return app, nil
}

func (a *LKSApp) initHttpServer(cfg *config.Config, log *zerolog.Logger) error {
	l := log.With().Str("context", "http_server").Logger()
	a.http = server.NewServer(a.db, &cfg.Http, &l)
	return nil
}

func (a *LKSApp) Run() error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		<-sigCh
		a.log.Info().Msg("catch signal to stop")
		if err := a.http.Stop(); err != nil {
			a.log.Error().Err(err).Msg("failed to stop http server")
		}
	}()

	a.log.Info().Msg("app start")
	if err := a.http.Run(); err != nil {
		a.log.Error().Err(err).Msg("failed to start http server")
		return errors.Wrap(err, "failed to start http server")
	}

	return nil
}

func (a *LKSApp) Stop() error {
	a.log.Info().Msg("app stop")
	return a.http.Stop()
}
