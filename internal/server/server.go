package server

import (
	"context"
	"github.com/daemondxx/lks_back/internal/config"
	"github.com/daemondxx/lks_back/internal/dao"
	ctrl_account "github.com/daemondxx/lks_back/internal/server/controllers/account"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	account_service "github.com/daemondxx/lks_back/internal/server/services/account"
	"github.com/daemondxx/lks_back/internal/server/services/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	cfg *config.HttpConfig
	eng *gin.Engine
	srv *http.Server
}

func NewServer(db *gorm.DB, cfg *config.HttpConfig, log *zerolog.Logger) *Server {
	r := gin.Default()
	r.Use(middleware.ErrorMiddleware())

	accDAO := dao.NewAccountDAO(db)
	tokenServ := token.NewTokenService(token.Config{Secret: []byte(cfg.JWTSecret)})
	accServ := account_service.NewAccountService(accDAO, tokenServ, log)
	accCtrl := ctrl_account.NewAccountController(accServ)

	{
		gr := r.Group("/account")
		accCtrl.ApplyAuthController(gr)
	}

	return &Server{
		eng: r,
		srv: &http.Server{
			Addr:    ":" + strconv.FormatUint(uint64(cfg.Port), 10),
			Handler: r.Handler(),
		},
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
