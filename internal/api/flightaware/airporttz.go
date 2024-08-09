package flightaware

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"time"
)

//go:embed airports.json
var tzData []byte

type tz struct {
	Code     string `json:"code"`
	Timezone string `json:"timezone"`
}

type airportTZ struct {
	arr []tz
}

func newAirportTZ() (*airportTZ, error) {
	dec := json.NewDecoder(bytes.NewReader(tzData))
	d := make([]tz, 0, 20698)
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}
	return &airportTZ{arr: d}, nil
}

func (t *airportTZ) GetLocation(airport string) (*time.Location, error) {
	if i, ok := slices.BinarySearchFunc(t.arr, tz{
		Code: airport,
	}, compare); ok {
		l, err := time.LoadLocation(t.arr[i].Timezone)
		if err != nil {
			return nil, fmt.Errorf("load location error: %w", err)
		}
		return l, nil
	} else {
		return nil, fmt.Errorf("not found tz in db")
	}
}

func compare(a tz, b tz) int {
	return cmp.Compare(a.Code, b.Code)
}
