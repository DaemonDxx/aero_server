package flightaware

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/daemondxx/lks_back/entity"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timeRegexp, _    = regexp.Compile("[0-9]{2}[:][0-9]{2}")
	offsetRegexp, _  = regexp.Compile("\\(([+-])\\d")
	airportRegexp, _ = regexp.Compile("[a-zA-Z]{3}")
	hourRegexp, _    = regexp.Compile("(\\d*)[ч, h]")
	minutesRegexp, _ = regexp.Compile(" (\\d*)[м, m]")
)

const url = "https://www.flightaware.com/live/flight/"

var (
	ErrFlightTimeExtraction = errors.New("flight time extraction from site error")
	ErrNotFoundRow          = fmt.Errorf("len rows with history data to be 0")
)

type flightHistoryRaw struct {
	Duration time.Duration
	Url      string
}

type bTab struct {
	ctx      context.Context
	id       int
	log      *zerolog.Logger
	rootLog  *zerolog.Logger
	tz       *airportTZ
	taskChan chan string
	resChan  chan *entity.FlightInfo
	errChan  chan error
}

func newTab(ctx context.Context, id int, tz *airportTZ, log *zerolog.Logger) (*bTab, error) {
	ctx, _ = chromedp.NewContext(ctx)
	if err := chromedp.Run(ctx); err != nil {
		return nil, fmt.Errorf("init main browser tab error: %e", err)
	}
	l := log.With().Int("tab_id", id).Logger()
	t := &bTab{
		ctx:      ctx,
		id:       id,
		tz:       tz,
		taskChan: make(chan string),
		resChan:  make(chan *entity.FlightInfo),
		errChan:  make(chan error),
		log:      &l,
		rootLog:  &l,
	}
	t.listen()

	return t, nil
}

func (t *bTab) listen() {
	go func() {
		for {
			flight := <-t.taskChan

			//goto actual flight page
			log := t.rootLog.With().Str("flight", flight).Logger()
			t.log = &log
			if err := t.gotoInfoPage(url + flight); err != nil {
				t.log.Err(err).Msg("navigate to page err")
				t.errChan <- err
				continue
			}

			rows, err := t.extractFlightHistoryRaw()

			if err != nil {
				t.log.Err(err).Msg("extract flight history rows error")
				t.errChan <- err
				continue
			}

			if len(rows) == 0 {
				t.log.Err(ErrNotFoundRow).Msg("not found dirty history row")
				t.errChan <- ErrNotFoundRow
				continue
			}

			var info *entity.FlightInfo
			info, err = t.extractTimeInfo()
			if err != nil {
				for _, r := range rows {
					err = nil
					if err := t.gotoInfoPage(r.Url); err != nil {
						t.log.Warn().Msg(fmt.Sprintf("go to %s error", r.Url))
						continue
					}
					info, err = t.extractTimeInfo()
					if err != nil {
						t.log.Warn().Msg(fmt.Sprintf("extract time info error: %e", err))
						continue
					}

					if info.Duration < 40*time.Minute {
						continue
					}
				}
			}

			if info == nil {
				t.errChan <- fmt.Errorf("cannot extract flight info for flight %s", flight)
				continue
			}

			rows = filterByDuration(rows)

			if len(rows) == 0 {
				t.log.Warn().Msg("len rows before filtering equal to be 0")
				info.AvgDuration = info.Duration
			} else {
				info.AvgDuration = average(rows)
			}

			info.FlightNumber = flight
			t.resChan <- info
		}
	}()
}

func (t *bTab) Execute(flight string) {
	t.taskChan <- flight
}

func (t *bTab) gotoInfoPage(url string) error {
	ctxTime, cancel := context.WithTimeout(t.ctx, 10*time.Second)
	defer cancel()

	if err := chromedp.Run(ctxTime,
		chromedp.Navigate(url),
	); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	ctxTime, cancel = context.WithTimeout(t.ctx, 30*time.Second)
	defer cancel()

	if err := chromedp.Run(ctxTime,
		chromedp.WaitVisible(".flightPageDataTableContainer", chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("wait selector error: %e", err)
	}
	return nil
}

func (t *bTab) extractTimeInfo() (*entity.FlightInfo, error) {
	ctx, cancel := context.WithTimeout(t.ctx, 5*time.Second)
	defer cancel()

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Nodes(".flightPageDataTimesChild > .flightPageDataAncillaryText > div > span", &nodes, chromedp.ByQueryAll),
	); err != nil {
		return nil, fmt.Errorf("extract nodes with time raw data.json error: %w", err)
	}

	var depTime, arvTime timeRaw
	for i, n := range nodes {
		if len(n.Children) == 0 {
			return nil, ErrFlightTimeExtraction
		}
		if i == 0 {
			if err := fillTimeRaw(n.Children[0].NodeValue, &depTime); err != nil {
				return nil, fmt.Errorf("fill dep time error: %w", err)
			}
		} else if i == 2 {
			if err := fillTimeRaw(n.Children[0].NodeValue, &arvTime); err != nil {
				return nil, ErrFlightTimeExtraction
			}
		}
	}

	if err := chromedp.Run(ctx,
		chromedp.Nodes(".flightPageSummaryAirportCode > .displayFlexElementContainer", &nodes, chromedp.ByQueryAll),
	); err != nil {
		return nil, fmt.Errorf("extract nodes with airport raw error: %w", err)
	}

	for i, n := range nodes[:2] {
		if len(n.Children) == 0 {
			return nil, fmt.Errorf("node child value with airport raw data not found")
		}
		airport := airportRegexp.FindString(n.Children[0].NodeValue)
		loc, err := t.tz.GetLocation(airport)
		if err != nil {
			return nil, fmt.Errorf("get location airport error: %w", err)
		}
		if i == 0 {
			depTime.Airport = airport
			depTime.Location = loc
		} else {
			arvTime.Airport = airport
			arvTime.Location = loc
		}
	}

	d := arvTime.sub(&depTime)

	return &entity.FlightInfo{
		From:          depTime.Airport,
		To:            arvTime.Airport,
		TimeDeparture: depTime.time(),
		Duration:      d,
	}, nil
}

func (t *bTab) extractFlightHistoryRaw() ([]flightHistoryRaw, error) {
	var nodes []*cdp.Node

	ctxTime, cancel := context.WithTimeout(t.ctx, 30*time.Second)
	defer cancel()

	if err := chromedp.Run(ctxTime,
		chromedp.Nodes(".flightPageDataRowTall", &nodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
	); err != nil || len(nodes) == 0 {
		return nil, fmt.Errorf("get flight history raws nodes error: %w", err)
	}

	rows := make([]flightHistoryRaw, 0)
	var child []*cdp.Node
	var c string
	var ok bool

	for i, n := range nodes {

		//Check is active row
		if c, _ = n.Attribute("class"); strings.Contains(c, "flightPageDataRowActive") {
			continue
		}
		c = ""

		//Extract url
		c, ok = n.Attribute("data-target")
		if !ok {
			t.log.Warn().Msg(fmt.Sprintf("cannot extract attribute 'data-target' from row %d", i))
			continue
		}

		//Extract Duration
		defer cancel()
		if err := chromedp.Run(ctxTime,
			chromedp.Nodes(".flightPageActivityLogData.optional.text-right > span", &child, chromedp.FromNode(n), chromedp.ByQuery, chromedp.AtLeast(0)),
		); err != nil || len(child) == 0 {
			t.log.Warn().Msg(fmt.Sprintf("extract duration time error (row %d): not found span", i))
			continue
		}

		if len(child[0].Children) == 0 {
			continue
		}

		sText := child[0].Children[0].NodeValue
		arr := hourRegexp.FindStringSubmatch(sText)
		if len(arr) == 0 {
			t.log.Warn().Msg(fmt.Sprintf("split hour (data %s) error (row %d)", sText, i))
			continue
		}
		h, err := strconv.Atoi(arr[1])
		if err != nil {
			t.log.Warn().Msg(fmt.Sprintf("parse hour (data %s) error (row %d)", arr[1], i))
			continue
		}
		arr = minutesRegexp.FindStringSubmatch(sText)
		m, err := strconv.Atoi(arr[1])
		if err != nil {
			t.log.Warn().Msg(fmt.Sprintf("parse minutes (data %s) error (row %d)", arr[1], i))
			continue
		}
		rows = append(rows, flightHistoryRaw{
			Duration: time.Duration(h*60+m) * time.Minute,
			Url:      "https://www.flightaware.com/" + c,
		})
		child = child[:0]
		c = ""
	}

	return rows, nil
}

func average(r []flightHistoryRaw) time.Duration {
	var sum uint64 = 0
	for _, row := range r {
		sum += uint64(row.Duration)
	}
	return time.Duration(sum / uint64(len(r)))
}
