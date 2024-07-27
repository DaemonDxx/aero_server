package lks

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/common"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

const tempNodeBuffSize = 10

var (
	ErrAccordAuth = errors.New("accord login or password incorrect")
	ErrLKSAuth    = errors.New("lks login or password incorrect")
)

var attrRegExp, _ = regexp.Compile("^text:\\s([^\\s]+)")

type AuthPayload struct {
	AccordLogin    string
	AccordPassword string
	LksLogin       string
	LksPassword    string
}

type LksAPIConfig struct {
}

type LksAPI struct {
	ctx context.Context
	log *zerolog.Logger

	//todo: оптимизировать выделение массивов под ноды
	tempNodes [][]*cdp.Node
}

func NewLksAPI(cfg *LksAPIConfig, log *zerolog.Logger) *LksAPI {
	l := log.With().Str("service", "lks_api").Logger()

	rootCtx, cancel := context.WithCancel(context.Background())
	ctx, _ := chromedp.NewExecAllocator(rootCtx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	ctx, _ = chromedp.NewContext(ctx)
	chromedp.Run(ctx)

	lks := &LksAPI{
		ctx: ctx,
		log: &l,
	}
	lks.initHandlerSignal(cancel)
	return lks
}

func (a *LksAPI) initHandlerSignal(c context.CancelFunc) {
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
	}()
}

func (a *LksAPI) acquireContext() context.Context {
	return a.ctx
}

func (a *LksAPI) releaseContext(ctx context.Context) {
	if err := chromedp.Run(ctx, common.ClearCookies()); err != nil {
		a.log.Err(err).Msg("clear cookies error")
	}
}

func (a *LksAPI) GetCurrentOrder(p AuthPayload) ([]CurrentOrderRow, error) {
	ctx := a.acquireContext()
	defer a.releaseContext(ctx)

	//_, err := a.GetMonthOrder(p)
	//return err
	if err := a.auth(ctx, p.AccordLogin, p.AccordPassword, p.LksLogin, p.LksPassword); err != nil {
		return nil, err
	}

	rows, err := a.extractOrder(ctx)
	if err != nil {
		return nil, err
	}
	return rows, nil

}

func (a *LksAPI) GetMonthOrder(p AuthPayload) (entity.Order, error) {
	var order entity.Order

	ctx := a.acquireContext()
	defer a.releaseContext(ctx)

	if err := a.auth(ctx, p.AccordLogin, p.AccordPassword, p.LksLogin, p.LksPassword); err != nil {
		return order, err
	}
	if err := a.extractMonthOrder(ctx); err != nil {
		return order, err
	}

	return order, nil
}

func (a *LksAPI) extractMonthOrder(ctx context.Context) error {
	log := a.getLogger("extract_month_order")

	ctx, _ = context.WithTimeout(ctx, 30*time.Second)
	var err error

	if err = chromedp.Run(ctx,
		chromedp.Navigate("https://lks.aeroflot.ru/AkkordOffice/PerspectivePlan"),
		common.WaitForNetworkIdle(),
	); err != nil {
		return err
	}

	var nodes []*cdp.Node
	if err = chromedp.Run(ctx,
		chromedp.Nodes("tbody[data-bind=\"foreach: perspectivePlan\"] > tr", &nodes, chromedp.ByQueryAll),
	); err != nil {
		return err
	}

	var childNodes []*cdp.Node
	for _, n := range nodes {
		if err = chromedp.Run(ctx,
			chromedp.Nodes("td", &childNodes, chromedp.ByQueryAll, chromedp.FromNode(n)),
		); err != nil {
			return err
		}

		for _, ch := range childNodes {
			log.Debug().Msg(ch.Children[0].NodeValue)
		}
	}

	return nil
}

func (a *LksAPI) auth(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) error {
	if err := a.authAccord(ctx, accLogin, accPass); err != nil {
		return err
	}
	if err := a.authLKS(ctx, lksLogin, lksPass); err != nil {
		return err
	}
	return nil
}

func (a *LksAPI) authAccord(ctx context.Context, login string, password string) error {
	log := a.getLogger("auth_accord").With().Str("accord_login", login).Logger()
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

func (a *LksAPI) authLKS(ctx context.Context, login string, password string) error {
	log := a.getLogger("auth_lks").With().Str("lks_login", login).Logger()
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

func (a *LksAPI) extractOrder(ctx context.Context) ([]CurrentOrderRow, error) {
	ctx, _ = context.WithTimeout(ctx, 30*time.Second)
	log := a.getLogger("extract_order")
	log.Debug().Msg("try extract order...")

	log.Debug().Msg("try navigate to current order site...")
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://lks.aeroflot.ru/AkkordOffice/CommitedRosters"),
		common.WaitForNetworkIdle(),
	); err != nil {
		log.Debug().Msgf("navigate error: %e", err)
		return nil, err
	}

	log.Debug().Msg("try extract nodes...")
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Nodes("tbody[data-bind=\"foreach: Duties\"] > tr", &nodes, chromedp.ByQueryAll),
	); err != nil {
		log.Debug().Msgf("extract nodes error: %s", err)
		return nil, err
	}

	var childNode []*cdp.Node
	rows := make([]CurrentOrderRow, 0, len(nodes))

	for i, n := range nodes {
		if err := chromedp.Run(ctx, chromedp.Nodes("td", &childNode, chromedp.ByQueryAll, chromedp.FromNode(n))); err != nil {
			return nil, err
		}
		m := extractMapFromDOM(childNode)
		row, err := newCOrderRow(m)
		if err != nil {
			log.Debug().Msg(fmt.Sprintf("extract row (%d) from lks error: %e", i, err))
			return rows, fmt.Errorf("extract row (%d) from lks error: %e", i, err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (a *LksAPI) getLogger(method string) zerolog.Logger {
	return a.log.With().Str("method", method).Logger()
}

func fillForm(formID string, inputs map[string]string) chromedp.Tasks {
	t := chromedp.Tasks{}
	for id, value := range inputs {
		t = append(t, chromedp.SetValue("#"+formID+" input[name='"+id+"']", value, chromedp.ByQuery))
	}
	t = append(t, chromedp.Submit("#"+formID, chromedp.ByQuery))
	return t
}

func extractMapFromDOM(nodes []*cdp.Node) map[string]string {
	r := map[string]string{}

	for _, n := range nodes {
		if attr, ok := n.Attribute("data-bind"); ok {
			if len(n.Children) != 0 {
				key := attrRegExp.FindStringSubmatch(attr)[1]
				r[key] = n.Children[0].NodeValue
			}
		}
	}
	return r
}
