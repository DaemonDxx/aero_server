package puller

import (
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/logger"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"time"
)

var defaultTime = time.Unix(1737786876, 0)

var suites = []struct {
	name          string
	pulledItems   []entity.OrderItem
	actualItems   []entity.OrderItem
	expectedID    []uint
	isHasErrEmpty bool
}{
	{
		name: "Empty actual, some pulled",
		pulledItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 10,
				},
				Departure: getTimeWithOffset(5),
			},
			{
				Model: gorm.Model{
					ID: 20,
				},
				Departure: getTimeWithOffset(7),
			},
		},
		actualItems:   []entity.OrderItem{},
		expectedID:    []uint{10, 20},
		isHasErrEmpty: false,
	},
	{
		name: "have actual, has not new pulled",
		pulledItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 10,
				},
				Departure: getTimeWithOffset(5),
			},
			{
				Model: gorm.Model{
					ID: 20,
				},
				Departure: getTimeWithOffset(7),
			},
		},
		actualItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Departure: getTimeWithOffset(5),
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Departure: getTimeWithOffset(7),
			},
		},
		expectedID:    []uint{},
		isHasErrEmpty: true,
	},
	{
		name: "have actual, has new pulled",
		pulledItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 20,
				},
				Departure: getTimeWithOffset(5),
			},
			{
				Model: gorm.Model{
					ID: 30,
				},
				Departure: getTimeWithOffset(6),
			},
			{
				Model: gorm.Model{
					ID: 40,
				},
				Departure: getTimeWithOffset(7),
			},
		},
		actualItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Departure: getTimeWithOffset(5),
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Departure: getTimeWithOffset(6),
			},
		},
		expectedID:    []uint{40},
		isHasErrEmpty: false,
	},
	{
		name: "have actual, has new pulled without crossing",
		pulledItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 20,
				},
				Departure: getTimeWithOffset(5),
			},
			{
				Model: gorm.Model{
					ID: 30,
				},
				Departure: getTimeWithOffset(6),
			},
		},
		actualItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Departure: getTimeWithOffset(3),
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Departure: getTimeWithOffset(4),
			},
		},
		expectedID:    []uint{20, 30},
		isHasErrEmpty: false,
	},
	{
		name:        "have actual, empty pulled",
		pulledItems: []entity.OrderItem{},
		actualItems: []entity.OrderItem{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Departure: getTimeWithOffset(3),
			},
		},
		expectedID:    []uint{},
		isHasErrEmpty: true,
	},
}

func TestDiffMethod(t *testing.T) {
	serv := Service{
		LoggedService: services.NewLoggedService("diff_test", logger.NewLogger("DEV")),
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			items, err := serv.diff(s.pulledItems, s.actualItems)
			if s.isHasErrEmpty {
				require.ErrorIs(t, err, errHasNotNewItems)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(s.expectedID), len(items))
				assert.Equal(t, s.expectedID, extractID(items))
			}
		})
	}
}

func extractID(i []entity.OrderItem) []uint {
	r := make([]uint, 0, len(i))
	for _, v := range i {
		r = append(r, v.ID)
	}
	return r
}

func getTimeWithOffset(o int) time.Time {
	return defaultTime.Add(time.Duration(o) * time.Hour)
}
