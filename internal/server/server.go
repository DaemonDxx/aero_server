package server

import (
	"context"
	"github.com/daemondxx/lks_back/internal/config"
	"github.com/daemondxx/lks_back/internal/dao"
	ctrl_account "github.com/daemondxx/lks_back/internal/server/controllers/account"
	ctrl_credential "github.com/daemondxx/lks_back/internal/server/controllers/credential"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	account_service "github.com/daemondxx/lks_back/internal/server/services/account"
	"github.com/daemondxx/lks_back/internal/server/services/authchecker"
	service_credential "github.com/daemondxx/lks_back/internal/server/services/credential"
	"github.com/daemondxx/lks_back/internal/server/services/token"
	lks_mock "github.com/daemondxx/lks_back/mocks/api/lks"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Server struct {
	cfg *config.Config
	log *zerolog.Logger
	eng *gin.Engine
	srv *http.Server
	db  *dao.Database
}

func NewServer(cfg *config.Config, log *zerolog.Logger) (*Server, error) {
	s := &Server{
		cfg: cfg,
		log: log,
	}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initHandlers(); err != nil {
		return nil, err
	}

	if err := s.initHttpServer(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) initDB() error {
	log.Info().Msg("init database...")

	d, err := dao.NewDatabase(s.cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("init database failed")
		return errors.Wrap(err, "failed to init database")
	}
	log.Info().Msg("database init successful")

	log.Info().Msg("start migrate database...")
	if err := d.AutoMigrate(); err != nil {
		log.Error().Err(err).Msg("failed to migrate database")
		return errors.Wrap(err, "failed to migrate database")
	}
	log.Info().Msg("migrate database successfully")

	s.db = d

	return nil
}

func (s *Server) initHandlers() error {
	db := s.db.GetInstance()

	r := gin.Default()
	r.Use(middleware.ErrorMiddleware())

	accDAO := dao.NewAccountDAO(db)
	credDAO := dao.NewCredentialDAO(db)

	tokenServ := token.NewTokenService(token.Config{Secret: []byte(s.cfg.Http.JWTSecret)})
	accServ := account_service.NewAccountService(accDAO, tokenServ, s.log)

	lksAPI := s.createLKSApi()
	checkerServ := authchecker.NewAuthCheckerService(lksAPI)
	credServ := service_credential.NewCredentialService(accServ, checkerServ, credDAO, s.log)

	auth := middleware.NewAuthMiddleware(tokenServ, accServ)

	accCtrl := ctrl_account.NewAccountController(accServ)
	credCtrl := ctrl_credential.NewCredentialController(credServ)

	{
		gr := r.Group("/account")
		accCtrl.ApplyAuthController(gr)
	}
	{
		gr := r.Group("/credential")
		gr.Use(auth.Handler)
		credCtrl.ApplyCredentialController(gr)
	}

	s.eng = r

	return nil
}

func (s *Server) createLKSApi() authchecker.LKSApi {
	switch s.cfg.Env {
	case "DEV":
		api := lks_mock.NewLKSApiMock(lks_mock.GetActualOrderFn_AccordLoginEquaslZero())
		return api
	}
	return nil
}

func (s *Server) initHttpServer() error {
	s.srv = &http.Server{
		Addr:    ":" + strconv.FormatUint(uint64(s.cfg.Http.Port), 10),
		Handler: s.eng.Handler(),
	}
	return nil
}

func (s *Server) Run() error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		<-sigCh
		s.log.Info().Msg("catch signal to stop")
		if err := s.Stop(); err != nil {
			s.log.Error().Err(err).Msg("failed to stop http server")
		}
	}()

	s.log.Info().Msg("server start listen")
	if err := s.srv.ListenAndServe(); err != nil {
		s.log.Error().Err(err).Msg("failed to start http server")
		return errors.Wrap(err, "failed to start http server")
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.log.Info().Msg("stop http server")
	return s.srv.Shutdown(ctx)
}
