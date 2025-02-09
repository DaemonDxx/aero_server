package collector

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
)

type AutoWorkerConfig struct {
	ActualOrderCronList []string
	MonthOrderCronList  []string
	TaskTimeout         time.Duration
}

type AutoCollectService struct {
	sch gocron.Scheduler
	log *zerolog.Logger
}

func NewAutoCollectService(s *Service, c *AutoWorkerConfig, log *zerolog.Logger) (*AutoCollectService, error) {
	l := log.With().Str("service", "auto_worker_collector").Logger()
	sch, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("create job scheduler error: %e", err)
	}

	for _, cOpt := range c.ActualOrderCronList {
		_, err := sch.NewJob(gocron.CronJob(cOpt, false), gocron.NewTask(func() {
			l.Info().Msg("start collect actual order...")
			ctx, cancel := context.WithTimeout(context.Background(), c.TaskTimeout)
			defer cancel()
			if err := s.CollectActualOrder(ctx); err != nil {
				var e *ErrLimitAttempt
				if errors.As(err, e) {
					for _, u := range e.Credentials {
						log.Warn().Msg(fmt.Sprintf("attempt limit for credential (id=%d) has been reached", u.ID))
					}
				} else {
					l.Err(err).Msg(fmt.Sprintf("collect actual orders error: %e", err))
				}
			} else {
				l.Info().Msg("collect actual order is successful")
			}
		}))
		if err != nil {
			return nil, fmt.Errorf("create actual order collect job error: %e", err)
		}
	}
	return &AutoCollectService{sch: sch, log: log}, nil
}

func (aw *AutoCollectService) Start() {
	aw.log.Info().Msg("auto worker collector start")
	aw.sch.Start()
}

func (aw *AutoCollectService) Stop() error {
	aw.log.Info().Msg("auto worker collector stop")
	return aw.sch.StopJobs()
}
