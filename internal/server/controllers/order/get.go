package order

import (
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/server/http_errors"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetOrdersResponse struct {
	Orders []entity.Order `json:"orders"`
}

func (ctrl *controller) Get(c *gin.Context) {
	acc, err := middleware.InjectAccountInfo(c)
	if err != nil {
		c.Error(err)
		return
	}

	if acc.CredentialID == nil {
		c.Error(http_errors.NewErrBadRequest("credential is not connected to this account"))
		return
	}

	o, err := ctrl.serv.GetActualOrders(c, *acc)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, GetOrdersResponse{
		Orders: o,
	})
}
