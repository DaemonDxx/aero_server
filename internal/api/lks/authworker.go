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
	"sync"
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
	rootCtx    context.Context
	ctx        context.Context
	tCh        chan authTask
	log        *zerolog.Logger
	opts       []chromedp.ExecAllocatorOption
	cancelFunc context.CancelFunc
	mu         sync.Mutex
}

func newWorker(ctx context.Context, tCh chan authTask, log *zerolog.Logger, debug bool) *worker {
	w := &worker{
		rootCtx: ctx,
		tCh:     tCh,
		log:     log,
		mu:      sync.Mutex{},
	}
	w.opts = append(w.opts, chromedp.Flag("headless", !debug))
	w.init()

	return w
}

func (w *worker) init() {
	if err := w.initCtx(); err != nil {
		w.log.Fatal().Msg(fmt.Sprintf("init browser context error: %e", err))
	}
	w.initOnCancelCtxHandle()
}

func (w *worker) initCtx() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	execCtx, cancel := chromedp.NewExecAllocator(w.rootCtx, append(chromedp.DefaultExecAllocatorOptions[:], w.opts...)...)
	w.cancelFunc = cancel

	ctx, _ := chromedp.NewContext(execCtx)
	w.ctx = ctx

	return chromedp.Run(ctx)
}

func (w *worker) initOnCancelCtxHandle() {
	go func() {
		for {
			<-w.ctx.Done()
			w.cancelFunc()

			w.log.Warn().Msg("browser context canceled. recreate context...")
			if err := w.initCtx(); err != nil {
				w.log.Fatal().Msg(fmt.Sprintf("context init error: %e", err))
			}
			w.log.Info().Msg("browser context is recreated")
		}
	}()
}

func (w *worker) listen() {
	go func() {
		for {
			t := <-w.tCh
			w.mu.Lock()
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
			w.mu.Unlock()
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
		log.Debug().Msgf("fill form error: %e", err)
		return err
	}

	log.Debug().Msg("check toast by auth error message...")
	var node []*cdp.Node
	if err := chromedp.Run(newCtx,
		chromedp.Nodes("#credentials_table_postheader > font", &node, chromedp.ByQueryAll, chromedp.AtLeast(0)),
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
		w := newWorker(rootCtx, taskCh, &wLog, cfg.Debug)
		w.listen()
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

func (a *AuthWorkerPool) Auth(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) (map[string]string, error) {
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
