package ctrl_credential

import (
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	service_credential "github.com/daemondxx/lks_back/internal/server/services/credential"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type CreateCredentialsBody struct {
	AccordLogin    string `json:"accordLogin" form:"accordLogin" binding:"required"`
	AccordPassword string `json:"accordPassword" form:"accordPassword" binding:"required"`
	LKSLogin       string `json:"lksLogin" form:"lksLogin" binding:"required"`
	LKSPassword    string `json:"lksPassword" form:"lksPassword" binding:"required"`
}

type CreateCredentialsResponse struct {
	ID          uint   `json:"id"`
	AccordLogin string `json:"accordLogin"`
	LksLogin    string `json:"lksLogin"`
}

func (ctrl *controller) CreateCredential(c *gin.Context) {
	var body CreateCredentialsBody

	acc, err := middleware.InjectAccountInfo(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	cr, err := ctrl.crServ.Create(c, acc, body.AccordLogin, body.AccordPassword, body.LKSLogin, body.LKSPassword)
	if err != nil {
		if errors.Is(err, lks.ErrLKSAuth) || errors.Is(err, lks.ErrAccordAuth) {
			var system string
			if errors.Is(err, lks.ErrAccordAuth) {
				system = "ACCORD"
			} else {
				system = "LKS"
			}

			c.JSON(http.StatusForbidden, gin.H{
				"message": "invalid login or password",
				"system":  system,
			})
			return
		} else if errors.Is(err, service_credential.ErrCredentialIsConnected) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "credential for this account is connected",
			})
			return
		} else {
			c.Error(err)
			return
		}
	}

	c.JSON(http.StatusCreated, CreateCredentialsResponse{
		ID:          cr.ID,
		AccordLogin: cr.AccordLogin,
		LksLogin:    cr.LKSLogin,
	},
	)
}
