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

type AutoWorker struct {
	sch gocron.Scheduler
	log *zerolog.Logger
}

func NewAutoWorker(s *Service, c *AutoWorkerConfig, log *zerolog.Logger) (*AutoWorker, error) {
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
				if errors.As(err, &ErrLimitAttempt{}) {
					for _, u := range err.(*ErrLimitAttempt).Users {
						log.Warn().Msg(fmt.Sprintf("attempt limit for user (id=%d) has been reached", u.ID))
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
	return &AutoWorker{sch: sch, log: log}, nil
}

func (aw *AutoWorker) Start() {
	aw.log.Info().Msg("auto worker collector start")
	aw.sch.Start()
}

func (aw *AutoWorker) Stop() error {
	aw.log.Info().Msg("auto worker collector stop")
	return aw.sch.StopJobs()
}
