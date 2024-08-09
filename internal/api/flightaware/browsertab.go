package flightaware

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/daemondxx/lks_back/entity"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"time"
)

type bTab struct {
	ctx      context.Context
	id       int
	log      *zerolog.Logger
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
	}
	t.listen()

	return t, nil
}

func (t *bTab) listen() {
	go func() {
		for {
			flight := <-t.taskChan
			//goto actual flight page
			log := t.log.With().Str("flight", flight).Logger()
			t.log = &log
			if err := t.gotoInfoPage(url + flight); err != nil {
				t.log.Err(err).Msg("navigate to page err")
				t.errChan <- err
				return
			}

			rows, err := t.extractFlightHistoryRaw()

			if err != nil {
				t.log.Err(err).Msg("extract flight history rows error")
				t.errChan <- err
				return
			}

			rows = filterByDuration(rows)

			var info *entity.FlightInfo
			info, err = t.extractTimeInfo()
			if err != nil {
				for _, r := range rows {
					err = nil
					if err := t.gotoInfoPage(r.Url); err != nil {
						t.log.Err(err).Msg(fmt.Sprintf("go to %s error", r.Url))
						continue
					}
					info, err = t.extractTimeInfo()
					if err != nil {
						t.log.Err(err).Msg("extract time info error")
					}
				}
			}

			if info == nil {
				t.errChan <- fmt.Errorf("cannot extract flight info for flight %s", flight)
				return
			}

			info.FlightNumber = flight
			info.AvgDuration = average(rows)
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

	ctxTime, cancel = context.WithTimeout(t.ctx, 60*time.Second)
	defer cancel()

	if err := chromedp.Run(ctxTime,
		chromedp.WaitVisible(".flightPageProgressTotal", chromedp.ByQuery),
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
		return nil, fmt.Errorf("extract nodes with airport raw data.json error: %w", err)
	}

	for i, n := range nodes[:2] {
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
		chromedp.Nodes(".flightPageDataRowTall", &nodes, chromedp.ByQueryAll),
	); err != nil {
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
		//Check is not unknown row
		if err := chromedp.Run(ctxTime,
			chromedp.Nodes(fmt.Sprintf(".flightPageDataRowTall:nth-child(%d) .flightPageResultUnknown", i+2), &child, chromedp.FromNode(n), chromedp.AtLeast(0)),
		); err != nil {
			t.log.Err(err).Msg("find unknowns row error")
			return nil, err
		} else if len(child) != 0 {
			child = child[:0]
			continue
		}

		//Extract url
		c, ok = n.Attribute("data-target")
		if !ok {
			t.log.Warn().Msg("cannot extract attribute 'data.json-target'")
			continue
		}

		//Extract Duration
		if err := chromedp.Run(ctxTime,
			chromedp.Nodes(".flightPageActivityLogData.optional.text-right > span", &child, chromedp.FromNode(n), chromedp.ByQuery),
		); err != nil {
			return nil, fmt.Errorf("extract duration node from row error: %e", err)
		} else if len(child) == 0 {
			return nil, fmt.Errorf("span with duration info not found")
		}

		arr := hourRegexp.FindStringSubmatch(child[0].Children[0].NodeValue)
		if len(arr) == 0 {
			continue
		}
		h, err := strconv.Atoi(arr[1])
		if err != nil {
			continue
		}
		arr = minutesRegexp.FindStringSubmatch(child[0].Children[0].NodeValue)
		m, err := strconv.Atoi(arr[1])
		if err != nil {
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
