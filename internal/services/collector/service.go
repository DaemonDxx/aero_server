package collector

import (
	"context"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"time"
)

const servName = "collector_service"
const defaultMaxAttempts = 1
const defaultMinTimeoutRetry = 1 * time.Minute
const defaultTimeoutContext = 30 * time.Second

type Config struct {
	MaxAttempts     uint
	MinTimeoutRetry time.Duration
}

type Service struct {
	services.LoggedService
	uDAO       UserDAO
	notifyServ NotificationService
	oServ      OrderService
	cfg        Config
}

func NewCollectorService(uDAO UserDAO, oServ OrderService, n NotificationService, cfg Config, log *zerolog.Logger) *Service {
	if log == nil {
		var l zerolog.Logger
		l = zerolog.New(os.Stdout).Level(zerolog.NoLevel)
		log = &l
	}

	if int64(cfg.MinTimeoutRetry) == 0 {
		cfg.MinTimeoutRetry = defaultMinTimeoutRetry
	}
	if cfg.MaxAttempts == 0 {
		cfg.MaxAttempts = defaultMaxAttempts
	}

	return &Service{
		LoggedService: services.NewLoggedService(log),
		uDAO:          uDAO,
		notifyServ:    n,
		oServ:         oServ,
		cfg:           cfg,
	}
}

func (s *Service) CollectActualOrder(ctx context.Context) error {
	log := s.GetLogger("collect_actual_order")

	log.Debug().Msg("find all active user")
	users, err := s.uDAO.Find(ctx, &entity.User{
		IsActive: true,
	})
	if err != nil {
		return &services.ErrServ{
			Service: servName,
			Message: "find all active user error",
			Err:     err,
		}
	}

	l := newUserList(users)

	attempt := 0
	var u *entity.User
	startTryTime := time.Now()

	for l.Len() != 0 {
		if uint(attempt) >= s.cfg.MaxAttempts {
			uFailed := l.Array()
			err := newErrLimitAttempt(uFailed)
			for _, usr := range uFailed {
				s.notifyServ.ErrorNotify(usr.ID, err)
			}
			return &services.ErrServ{
				Service: servName,
				Message: err.Error(),
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
			u = el.v

			c, cancel := context.WithTimeout(ctx, defaultTimeoutContext)
			o, err := s.oServ.GetActualOrder(c, u.ID)
			cancel()

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					if el.next == nil {
						break
					} else {
						el = el.next
						continue
					}
				}

				log.Err(err).Msg(fmt.Sprintf("collect actual order for user (id=%d) error: %e", u.ID, err))
				s.notifyServ.ErrorNotify(u.ID, err)
			} else {
				s.notifyServ.ActualOrderNotify(u.ID, o)
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
