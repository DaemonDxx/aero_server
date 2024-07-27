package lks

import (
	"strings"
	"time"
)

func parseLKSFormatData(s string) (time.Time, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("02.01.06 15:04", strings.Trim(s, " "), loc)
}

func parseLKSConfirmedFormatData(s string) (time.Time, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("02.01.2006 15:04", strings.Trim(s, " "), loc)
}
