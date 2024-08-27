package order

import (
	"github.com/daemondxx/lks_back/internal/api/lks"
	"time"
)

func getFirstDuty(confirmDate *time.Time) []lks.CurrentDuty {
	return []lks.CurrentDuty{
		{
			Code:         "",
			FlightNumber: "1048 2928 2929 1049",
			AircraftType: "320",
			Route:        "Ш (B) Сочи Омск Сочи Ш",
			StartDate:    time.Date(2024, 8, 8, 16, 15, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 11, 14, 25, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Явка в терминал B, 2 этаж, комната брифинга  2.12.007",
			BlockDate:    nil,
		},
		{
			Code:         "",
			FlightNumber: "2134 2135",
			AircraftType: "321",
			Route:        "Ш (C) Стамбул Ш",
			StartDate:    time.Date(2024, 8, 12, 14, 15, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 13, 01, 00, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Явка в терминал C",
			BlockDate:    nil,
		},
	}
}

func getSecondDuty(confirmDate *time.Time) []lks.CurrentDuty {
	return []lks.CurrentDuty{
		{
			Code:         "",
			FlightNumber: "1230 1231",
			AircraftType: "321",
			Route:        "Ш (B) Уфа Ш",
			StartDate:    time.Date(2024, 8, 17, 6, 20, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 17, 11, 40, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Явка в терминал B, 2 этаж, комната брифинга  2.12.007",
			BlockDate:    nil,
		},
		{
			Code:         "",
			FlightNumber: "1124 1125",
			AircraftType: "321",
			Route:        "Ш (B) Сочи Ш",
			StartDate:    time.Date(2024, 8, 19, 9, 55, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 13, 17, 50, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Явка в терминал B, 2 этаж, комната брифинга  2.12.007",
			BlockDate:    nil,
		},
		{
			Code:         "",
			FlightNumber: "1280 1281",
			AircraftType: "320",
			Route:        "Ш (B) Казань Ш",
			StartDate:    time.Date(2024, 8, 20, 17, 00, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 13, 21, 30, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Явка в терминал B, 2 этаж, комната брифинга  2.12.007",
			BlockDate:    nil,
		},
	}
}

func getNoFlightsDuty(confirmDate *time.Time) []lks.CurrentDuty {
	return []lks.CurrentDuty{
		{
			Code:         "",
			FlightNumber: "",
			AircraftType: "",
			Route:        "",
			StartDate:    time.Date(2024, 8, 17, 6, 20, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 17, 11, 40, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Резер домашний дневной",
			BlockDate:    nil,
		},
		{
			Code:         "",
			FlightNumber: "",
			AircraftType: "",
			Route:        "",
			StartDate:    time.Date(2024, 8, 17, 6, 20, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 17, 11, 40, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Резер домашний ночной",
			BlockDate:    nil,
		},
		{
			Code:         "",
			FlightNumber: "",
			AircraftType: "",
			Route:        "",
			StartDate:    time.Date(2024, 8, 17, 6, 20, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 17, 11, 40, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Отпуск",
			BlockDate:    nil,
		},
		{
			Code:         "",
			FlightNumber: "",
			AircraftType: "",
			Route:        "",
			StartDate:    time.Date(2024, 8, 17, 6, 20, 0, 0, &time.Location{}),
			EndDate:      time.Date(2024, 8, 17, 11, 40, 0, 0, &time.Location{}),
			ConfirmDate:  confirmDate,
			ConfirmType:  4,
			Target:       "",
			Place:        "",
			Note:         "Other",
			BlockDate:    nil,
		},
	}
}
