package collector

import (
	"context"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/server/services/collector/payloads"
	"github.com/daemondxx/lks_back/internal/server/services/notifier"
	order "github.com/daemondxx/lks_back/internal/server/services/order"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/pkg/errors"
	"time"
)

const defaultMaxAttempts = 1
const defaultMinTimeoutRetry = 1 * time.Minute
const defaultTimeoutContext = 30 * time.Second

func (s *Service) CollectActualOrder(ctx context.Context, crs []entity.Credential) error {
	log := s.GetLogger("collect_actual_order")

	var cr *entity.Credential
	attempt := 0
	startTryTime := time.Now()
	l := newCredentialList(crs)

	for l.Len() != 0 {
		if attempt >= s.cfg.MaxAttempts {
			credFailed := l.Array()
			err := newErrLimitAttempt(credFailed)

			go func(crs []*entity.Credential) {
				p := payloads.NewInternalErrorPayload()
				for _, cr := range crs {
					s.notify.Send(notifier.Message{
						To:      *cr,
						Payload: p,
					})
				}
			}(credFailed)

			return &services.ErrServ{
				Service: servName,
				Message: "max attempts reached",
				Err:     err,
			}
		}

		log.Info().Msg(fmt.Sprintf("start %d attempt collect orders", attempt+1))
		if attempt > 0 {
			d := time.Now().Sub(startTryTime)
			if d < s.cfg.MinTimeoutRetry {
				log.Info().Msg(fmt.Sprintf("the timeout between attempts has not expired. wait %d ms", (s.cfg.MinTimeoutRetry-d)/time.Millisecond))
				t := time.NewTimer(s.cfg.MinTimeoutRetry - d)
				<-t.C
				log.Info().Msg("timeout expired. continue collect")
				startTryTime = time.Now()
			}
		}

		el := l.First()

		for {
			cr = el.v

			c, cancel := context.WithTimeout(ctx, defaultTimeoutContext)
			o, err := s.oServ.PullNewOrder(c, cr)
			cancel()

			if err != nil {
				if errors.Is(err, order.ErrEmptyOrder) {
					_, err := s.oServ.Create(ctx, cr, nil)
					if err != nil {
						s.notify.Send(notifier.Message{
							To:      *cr,
							Payload: payloads.NewInternalErrorPayload(),
						})
						log.Err(err).Uint("credential", cr.ID).Msg("create empty order failed")
					} else {
						s.notify.Send(notifier.Message{
							To:      *cr,
							Payload: payloads.NewUpdateOrderPayload(o),
						})
					}
				} else if errors.Is(err, lks.ErrLKSAuth) || errors.Is(err, lks.ErrAccordAuth) {
					if err := s.crDAO.SetInactiveStatus(ctx, cr); err != nil {
						log.Err(err).Uint("credential", cr.ID).Msg("set inactive status failed")
					}
					s.notify.Send(notifier.Message{
						To:      *cr,
						Payload: payloads.NewAuthErrorPayload(),
					})
				} else if errors.Is(err, context.DeadlineExceeded) {
					if el.next == nil {
						break
					} else {
						el = el.next
						continue
					}
				} else {
					log.Err(err).Msg(fmt.Sprintf("collect actual order for user (id=%d) error: %e", cr.ID, err))
					s.notify.Send(notifier.Message{
						To:      *cr,
						Payload: payloads.NewInternalErrorPayload(),
					})
				}
			}

			n := el.next
			l.Remove(el)

			if n == nil {
				break
			} else {
				el = n
			}
		}
		attempt++
	}

	return nil
}
