package lks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var cases = []struct {
	fixture string
	err     bool
	name    string
	outUTC  int64
}{
	{
		fixture: "28.06.24 12:00",
		err:     false,
		name:    "simple parse test ",
		outUTC:  1719565200000,
	},
	{
		fixture: " 28.06.24 12:00  ",
		err:     false,
		name:    "parse date with around spaces",
		outUTC:  1719565200000,
	},
	{
		fixture: "28.06. 12:00",
		err:     true,
		name:    "parse date without year",
		outUTC:  0,
	},
	{
		fixture: "",
		err:     true,
		name:    "parse date from empty string",
		outUTC:  0,
	},
}

func TestParseDate(t *testing.T) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			time, err := parseLKSFormatData(tc.fixture)
			if !tc.err {
				assert.NoError(t, err)
				assert.Equal(t, tc.outUTC, time.UTC().UnixMilli())
			} else {
				assert.Error(t, err)
			}

		})
	}
}
