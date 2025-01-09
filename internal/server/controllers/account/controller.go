package ctrl_account

import (
	"context"
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	CreateAccount(ctx context.Context, tgID uint64) (string, error)
}

type controller struct {
	serv AccountService
}

func NewAccountController(serv AccountService) *controller {
	return &controller{
		serv: serv,
	}
}

func (ctrl *controller) ApplyAuthController(r *gin.RouterGroup) {
	r.POST("/", ctrl.CreateAccount)
}
