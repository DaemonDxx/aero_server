package flightaware

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
	"sync"
)

type initBrowserFunc func(b *browser)

var withLogger = func(log *zerolog.Logger) initBrowserFunc {
	return func(b *browser) {
		b.log = log
	}
}

var withDebugMode = func() initBrowserFunc {
	return func(b *browser) {
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
		b.ctx = ctx
		b.cancelFunc = cancel
	}
}

var withProductionMode = func() initBrowserFunc {
	return func(b *browser) {
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", true))...)
		b.ctx = ctx
		b.cancelFunc = cancel
	}
}

type browser struct {
	ctx           context.Context
	cancelFunc    context.CancelFunc
	tabs          []*bTab
	availableTabs chan *bTab
	mu            sync.RWMutex
	log           *zerolog.Logger
	tz            *airportTZ
}

func newBrowser(maxTabCount uint, initParams ...initBrowserFunc) (*browser, error) {
	b := browser{
		mu: sync.RWMutex{},
	}

	for _, f := range initParams {
		f(&b)
	}

	if b.log == nil {
		l := zerolog.New(zerolog.ConsoleWriter{
			Out: nil,
		}).Level(zerolog.Disabled)
		b.log = &l
	}

	if b.ctx == nil {
		withProductionMode()(&b)
	}

	b.log.Debug().Msg("init tz database...")
	tz, err := newAirportTZ()
	if err != nil {
		b.log.Err(err).Msg("database init error")
		return nil, fmt.Errorf("tz database init error: %e", err)
	}
	b.log.Debug().Msg("tz database init successful")
	b.tz = tz

	b.log.Debug().Msg("init browser tabs...")
	if err := b.initTabs(maxTabCount); err != nil {
		b.log.Err(err).Msg("init tabs browser error")
		return nil, err
	}
	b.log.Debug().Msg("browser tabs init successful")

	return &b, nil
}

func (b *browser) initTabs(countTabs uint) error {
	b.tabs = make([]*bTab, countTabs)
	b.availableTabs = make(chan *bTab, countTabs)

	for i := 0; i < int(countTabs); i++ {
		b.log.Debug().Msg(fmt.Sprintf("try init tab №%d", i))
		t, err := newTab(b.ctx, i, b.tz, b.log)
		if err != nil {
			return err
		}
		b.log.Debug().Msg(fmt.Sprintf("tab №%d init successful", i))

		b.tabs[i] = t
		b.availableTabs <- t

		if i == 0 {
			b.ctx = b.tabs[0].ctx
		}
	}

	return nil
}

func (b *browser) GetFlightInfo(ctx context.Context, flight string) (Information, error) {
	t := <-b.availableTabs

	t.Execute(flight)

	select {
	case info := <-t.resChan:
		b.availableTabs <- t
		return *info, nil
	case err := <-t.errChan:
		//todo подумать над пересозданием вкладки
		b.availableTabs <- t
		return Information{}, err
	case <-ctx.Done():
		go func() {
			select {
			case <-t.resChan:
			case <-t.errChan:
				b.availableTabs <- t
			}
		}()
		return Information{}, ctx.Err()
	}
}

func (b *browser) Close() error {
	b.cancelFunc()
	return nil
}
