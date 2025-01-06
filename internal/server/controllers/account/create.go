package ctrl_account

import (
	httpErr "github.com/daemondxx/lks_back/internal/server/http_errors"
	service_account "github.com/daemondxx/lks_back/internal/server/services/account"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type CreateAccountBody struct {
	//todo добавить валидацию
	TelegramID uint64 `json:"telegramID" form:"telegramID" binding:"required"`
}

type CreateAccountResponse struct {
	AccessToken string `json:"accessToken"`
}

// CreateAccount godoc
//
//	@Summary		Create account
//	@Description	Create new account by telegram id
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Param			b	body		account.CreateAccountBody	true	"Create account DTO"
//	@Success		201	{object}	account.CreateAccountResponse
//	@Failure		400	{object}	http_errors.HttpError
//	@Failure		500	{object}	http_errors.HttpError
//	@Router			/account/ [post]
func (ctrl *controller) CreateAccount(c *gin.Context) {
	var body CreateAccountBody

	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	t, err := ctrl.serv.CreateAccount(c, body.TelegramID)
	if err != nil {
		if errors.Is(err, service_account.ErrAccountExists) {
			c.Error(httpErr.NewErrBadRequest("account already exist"))
		} else {
			c.Error(err)
		}
		return
	}

	c.JSON(http.StatusCreated, &CreateAccountResponse{
		AccessToken: t,
	})
}
