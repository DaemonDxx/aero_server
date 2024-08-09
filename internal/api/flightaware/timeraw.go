package flightaware

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type timeRaw struct {
	Hour      int
	Minutes   int
	Airport   string
	Location  *time.Location
	OffsetDay int
}

func fillTimeRaw(s string, tr *timeRaw) error {
	t := timeRegexp.FindString(s)
	timeRes := strings.Split(t, ":")
	h, err := strconv.Atoi(timeRes[0])
	if err != nil {
		return fmt.Errorf("parse raw hour error: %s", timeRes[0])
	}
	m, err := strconv.Atoi(timeRes[1])
	if err != nil {
		return fmt.Errorf("parse raw minutes error: %s", timeRes[1])
	}
	tr.Hour = h
	tr.Minutes = m

	offsetRes := offsetRegexp.FindStringSubmatch(s)
	if len(offsetRes) == 2 {
		if offsetRes[1] == "+" {
			tr.OffsetDay = 1
		} else {
			tr.OffsetDay = -1
		}
	}
	return nil
}

func (tr *timeRaw) sub(t *timeRaw) time.Duration {
	return tr.time().Sub(t.time())
}

func (tr *timeRaw) time() time.Time {
	c := time.Now()
	return time.Date(c.Year(), c.Month(), c.Day()+tr.OffsetDay, tr.Hour, tr.Minutes, 0, 0, tr.Location)
}
