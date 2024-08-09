package flightaware

import (
	"context"
	_ "embed"
	"errors"
	"github.com/daemondxx/lks_back/entity"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

const url = "https://www.flightaware.com/live/flight/"

var (
	timeRegexp, _    = regexp.Compile("[0-9]{2}[:][0-9]{2}")
	offsetRegexp, _  = regexp.Compile("\\(([+-])\\d")
	airportRegexp, _ = regexp.Compile("[a-zA-Z]{3}")
	hourRegexp, _    = regexp.Compile("(\\d*)[ч]")
	minutesRegexp, _ = regexp.Compile("(\\d*)[м]")
)

var (
	ErrFlightTimeExtraction = errors.New("flight time extraction from site error")
)

type flightHistoryRaw struct {
	Duration time.Duration
	Url      string
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
	l := log.With().Str("service", "flightaware_api").Logger()
	b, err := newBrowser(uint(cfg.MaxTabCount), withLogger(&l), withDebugMode())
	if err != nil {
		return nil, err
	}

	api := &Api{
		browser: b,
		log:     &l,
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

func (f *Api) GetFlightInfo(ctx context.Context, flight string) (entity.FlightInfo, error) {
	return f.browser.GetFlightInfo(ctx, flight)
}
