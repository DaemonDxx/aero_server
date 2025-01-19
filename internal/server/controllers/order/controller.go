package order

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/gin-gonic/gin"
)

type OrderService interface {
	GetActualOrders(ctx context.Context, acc entity.Account) ([]entity.Order, error)
	ConfirmOrder(ctx context.Context, acc entity.Account, orderID uint) error
}

type controller struct {
	serv OrderService
}

func NewOrderController(serv OrderService) *controller {
	return &controller{
		serv: serv,
	}
}

func (ctrl *controller) ApplyHandlers(g *gin.RouterGroup) {
	g.GET("/", ctrl.Get)
}
