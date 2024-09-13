package flightaware

import (
	"context"
	_ "embed"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Information struct {
	FlightNumber  string
	From          string
	To            string
	TimeDeparture time.Time
	Duration      time.Duration
	AvgDuration   time.Duration
}

type ApiConfig struct {
	MaxTabCount int
	Debug       bool
}

type Api struct {
	browser *browser
	log     *zerolog.Logger
}

func NewFlightInfoAPI(cfg *ApiConfig, log *zerolog.Logger) (*Api, error) {
	if log == nil {
		var l zerolog.Logger
		l = zerolog.New(os.Stdout).Level(zerolog.NoLevel)
		log = &l
	}

	var bMode initBrowserFunc

	if cfg.Debug {
		bMode = withDebugMode()
	} else {
		bMode = withProductionMode()
	}

	b, err := newBrowser(uint(cfg.MaxTabCount), withLogger(log), bMode)
	if err != nil {
		return nil, err
	}

	api := &Api{
		browser: b,
		log:     log,
	}

	api.initHandlerSignal(b)

	return api, nil
}

func (f *Api) initHandlerSignal(b *browser) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		<-sigCh
		b.Close()
	}()
}

func (f *Api) GetFlightInfo(ctx context.Context, flight string) (Information, error) {
	return f.browser.GetFlightInfo(ctx, flight)
}
