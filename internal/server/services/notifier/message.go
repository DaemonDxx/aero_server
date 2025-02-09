package notifier

import (
	"fmt"
	"github.com/daemondxx/lks_back/entity"
)

type Message struct {
	To      entity.Credential
	Payload fmt.Stringer
}
