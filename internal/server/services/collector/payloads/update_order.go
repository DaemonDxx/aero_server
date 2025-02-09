package payloads

import (
	"fmt"
	"github.com/daemondxx/lks_back/entity"
)

type UpdateOrder struct {
	Order *entity.Order
}

func NewUpdateOrderPayload(o *entity.Order) *UpdateOrder {
	return &UpdateOrder{
		Order: o,
	}
}

func (p *UpdateOrder) String() string {
	if len(p.Order.Items) == 0 {
		return "Был получен пустой наряд"
	} else {
		return fmt.Sprintf("У вас %d новых рейса! Зайдите в приложение, чтобы узнать подробности!", len(p.Order.Items))
	}
}
