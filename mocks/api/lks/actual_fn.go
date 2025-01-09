package lks_mock

import (
	"context"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"strconv"
	"time"
)

func GetActualOrderFn_AccordLoginEquaslZero() ActualFn {
	return func(ctx context.Context, p lks.AuthPayload) ([]lks.CurrentDuty, error) {
		n, err := strconv.ParseInt(p.AccordLogin, 10, 64)

		if err != nil {
			return []lks.CurrentDuty{
				{
					Code:         "",
					FlightNumber: "1048 2928 2929 1049",
					AircraftType: "320",
					Route:        "Ш (B) Сочи Омск Сочи Ш",
					StartDate:    time.Date(2024, 8, 8, 16, 15, 0, 0, &time.Location{}),
					EndDate:      time.Date(2024, 8, 11, 14, 25, 0, 0, &time.Location{}),
					ConfirmDate:  nil,
					ConfirmType:  4,
					Target:       "",
					Place:        "",
					Note:         "Явка в терминал B, 2 этаж, комната брифинга  2.12.007",
					BlockDate:    nil,
				}}, nil
		}

		if n > 0 {
			return nil, lks.ErrAccordAuth
		} else {
			return nil, lks.ErrLKSAuth
		}
	}
}
