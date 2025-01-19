package ctrl_credential

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/gin-gonic/gin"
)

type CredentialService interface {
	Create(ctx context.Context, acc *entity.Account, lksLogin string, lksPass string, accLogin string, accPass string) (*entity.Credential, error)
	Update(ctx context.Context, id uint, accPass string, lksPass string) error
}

type controller struct {
	crServ CredentialService
}

func NewCredentialController(crServ CredentialService) *controller {
	return &controller{
		crServ: crServ,
	}
}

func (ctrl *controller) ApplyHandlers(r *gin.RouterGroup) {
	r.POST("/", ctrl.CreateCredential)
	r.PATCH("/:id", ctrl.UpdateCredential)
}
