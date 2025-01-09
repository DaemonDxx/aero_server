package ctrl_credential

import (
	"github.com/daemondxx/lks_back/internal/api/lks"
	httpErr "github.com/daemondxx/lks_back/internal/server/http_errors"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type UpdateCredentialBody struct {
	AccordPassword string `json:"accordPassword" form:"accordPassword" validate:"required"`
	LKSPassword    string `json:"lksPassword" form:"lksPassword" validate:"required"`
}

func (ctrl *controller) UpdateCredential(c *gin.Context) {
	var body UpdateCredentialBody

	acc, err := middleware.InjectAccountInfo(c)
	if err != nil {
		c.Error(err)
		return
	}

	p := c.Param("id")
	if p == "" {
		c.Error(httpErr.NewErrBadRequest("credential id param is required"))
		return
	}
	crID, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.Error(httpErr.NewErrBadRequest("credential id param is invalid"))
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if *acc.CredentialID != uint(crID) {
		c.Error(httpErr.NewErrForbidden("credential doesn't belong to you"))
		return
	}

	if err := ctrl.crServ.Update(c, uint(crID), body.AccordPassword, body.LKSPassword); err != nil {
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
		} else {
			c.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}
