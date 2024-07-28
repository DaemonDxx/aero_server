package lks

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/daemondxx/lks_back/internal/api/common"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type authTask struct {
	accLogin string
	accPass  string
	lksLogin string
	lksPass  string
	resCh    chan *taskResult
}

type taskResult struct {
	cookie map[string]string
	err    error
}

type worker struct {
	ctx context.Context
	tCh chan authTask
	log *zerolog.Logger
}

func newWorker(ctx context.Context, tCh chan authTask, log *zerolog.Logger, debug bool) *worker {
	w := &worker{
		ctx: ctx,
		tCh: tCh,
		log: log,
	}
	w.init(debug)
	return w
}

func (w *worker) init(debug bool) {
	var ctx context.Context

	if debug {
		ctx, _ = chromedp.NewExecAllocator(w.ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	} else {
		ctx, _ = chromedp.NewExecAllocator(w.ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", true))...)
	}

	ctx, _ = chromedp.NewContext(ctx)
	w.ctx = ctx
	chromedp.Run(ctx)

	go func() {
		for {
			t := <-w.tCh
			c, err := w.auth(w.ctx, t.accLogin, t.accPass, t.lksLogin, t.lksPass)
			//todo оптимазция создания результатов
			t.resCh <- &taskResult{
				cookie: c,
				err:    err,
			}
			//todo закрытие канала не в удачном месте
			close(t.resCh)
			if err := chromedp.Run(w.ctx, common.ClearCookies()); err != nil {
				w.log.Err(err).Msg("clear cookies error")
			}
		}
	}()
}

func (w *worker) auth(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) (map[string]string, error) {
	cookies := make(map[string]string)
	if err := w.authAccord(ctx, accLogin, accPass); err != nil {
		return cookies, err
	}
	if err := w.authLKS(ctx, lksLogin, lksPass); err != nil {
		return cookies, err
	}

	if err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			c, err := network.GetCookies().Do(ctx)
			if err != nil {
				fmt.Println(err)
				return err
			}
			for _, i := range c {
				cookies[i.Name] = i.Value
			}
			return nil
		}),
	); err != nil {
		return cookies, err
	}
	return cookies, nil
}

func (w *worker) authAccord(ctx context.Context, login string, password string) error {
	log := w.getLogger("auth_accord").With().Str("accord_login", login).Logger()
	log.Debug().Msg("try auth accord system")

	newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Debug().Msg("try navigate and wait form...")
	if err := chromedp.Run(newCtx,
		chromedp.Navigate("https://lks.aeroflot.ru/"),
		chromedp.WaitVisible("auth_form", chromedp.ByID),
	); err != nil {
		log.Debug().Msgf("try navigate error: %e", err)
		return err
	}

	newCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Debug().Msg("try fill form...")
	if err := chromedp.Run(newCtx,
		fillForm("auth_form", map[string]string{
			"username": login,
			"password": password,
		}),
		common.WaitForNetworkIdle(),
	); err != nil {
		fmt.Println(err)
		log.Debug().Msgf("fill form error: %e", err)
		return err
	}

	newCtx, cancel = context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	log.Debug().Msg("check toast by auth error message...")
	var node []*cdp.Node
	if err := chromedp.Run(newCtx,
		chromedp.Nodes("#credentials_table_postheader > font", &node, chromedp.ByQuery),
	); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			log.Debug().Msgf("check toast by auth error: %e", err)
			return err
		}
	}

	if len(node) != 0 {
		log.Debug().Msg("accord auth failed cause invalid arguments")
		return ErrAccordAuth
	}

	log.Debug().Msg("try auth successful")
	return nil
}

func (w *worker) authLKS(ctx context.Context, login string, password string) error {
	log := w.getLogger("auth_lks").With().Str("lks_login", login).Logger()
	log.Debug().Msg("try auth lks system...")

	ctxTime, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Debug().Msg("try fill form...")
	if err := chromedp.Run(ctxTime,
		chromedp.SetValue("#loginForm_Username", login, chromedp.ByQuery),
		chromedp.SetValue("#loginForm_Password", password, chromedp.ByQuery),
		chromedp.Click(".btn-primary", chromedp.ByQuery),
	); err != nil {
		log.Debug().Msgf("fill lks form error: %e", err)
		return fmt.Errorf("fill lks auth form error: %w", err)
	}

	ctxTime, cancel = context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	log.Debug().Msg("wait toast auth error")
	if err := chromedp.Run(ctxTime,
		chromedp.WaitVisible(".toast", chromedp.ByQuery),
	); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			log.Debug().Msgf("wait toast error: %e", err)
			return err
		}
	} else {
		log.Debug().Msg("lks auth failed cause invalid arguments")
		return ErrLKSAuth
	}

	ctxTime, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Debug().Msg("wait navigate to lks system...")
	if err := chromedp.Run(ctxTime,
		common.WaitForNetworkIdle(),
	); err != nil {
		log.Debug().Msgf("navigate to lks system error: %e", err)
		return err

	}
	log.Debug().Msg("lks auth successful")
	return nil
}

func (w *worker) getLogger(method string) zerolog.Logger {
	return w.log.With().Str("method", method).Logger()
}

type AuthWorkerPoolConfig struct {
	PoolSize uint
	Debug    bool
}

type AuthWorkerPool struct {
	ctx    context.Context
	taskCh chan authTask
}

func NewAuthWorkerPool(cfg *AuthWorkerPoolConfig, log *zerolog.Logger) *AuthWorkerPool {
	taskCh := make(chan authTask, cfg.PoolSize)
	rootCtx, cancel := context.WithCancel(context.Background())

	for i := 0; i < int(cfg.PoolSize); i++ {
		wLog := log.With().Str("service", "auth_worker_"+string(rune(i))).Logger()
		newWorker(rootCtx, taskCh, &wLog, cfg.Debug)
	}

	pool := &AuthWorkerPool{
		taskCh: taskCh,
	}
	pool.initHandlerSignal(cancel)
	return pool
}

func (a *AuthWorkerPool) initHandlerSignal(c context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		<-sigCh
		c()
		close(a.taskCh)
	}()
}

func (a *AuthWorkerPool) auth(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) (map[string]string, error) {
	//todo оптимизация получения канала
	resCh := make(chan *taskResult)

	t := authTask{
		accLogin: accLogin,
		accPass:  accPass,
		lksLogin: lksLogin,
		lksPass:  lksPass,
		resCh:    resCh,
	}
	a.taskCh <- t

	select {
	case res := <-resCh:
		return res.cookie, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

func fillForm(formID string, inputs map[string]string) chromedp.Tasks {
	t := chromedp.Tasks{}
	for id, value := range inputs {
		t = append(t, chromedp.SetValue("#"+formID+" input[name='"+id+"']", value, chromedp.ByQuery))
	}
	t = append(t, chromedp.Submit("#"+formID, chromedp.ByQuery))
	return t
}
