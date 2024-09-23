package implementation

import (
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	notifier "github.com/daemondxx/lks_back/internal/services/notifier"
	"github.com/pkg/errors"
)

const updateActualOrderKey = "UPDATE_ACTUAL_ORDER"
const unauthorizedErrorKey = "UNAUTHORIZED"
const collectInternalErrorKey = "COLLECT_INT_ERROR"

type CollectorNotifier struct {
	n *notifier.Service
}

func NewCollectorNotifier(n *notifier.Service) *CollectorNotifier {
	return &CollectorNotifier{
		n: n,
	}
}

func (c *CollectorNotifier) ActualOrderNotify(userID uint, o entity.Order) {
	n := notifier.Notification{
		FromUserID: userID,
		Key:        updateActualOrderKey,
		Payload: struct {
			Order entity.Order `json:"order"`
		}{
			Order: o,
		},
	}
	c.n.Notify(n)
}

func (c *CollectorNotifier) ErrorNotify(userID uint, err error) {
	n := notifier.Notification{
		FromUserID: userID,
	}

	if errors.Is(err, lks.ErrAccordAuth) || errors.Is(err, lks.ErrLKSAuth) {
		n.Key = unauthorizedErrorKey
	} else {
		n.Key = collectInternalErrorKey
	}

	c.n.Notify(n)
}
