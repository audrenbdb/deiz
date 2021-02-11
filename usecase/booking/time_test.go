package booking_test

import (
	"github.com/audrenbdb/deiz/usecase/booking"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeRangesOverlaps(t *testing.T) {
	var tests = []struct {
		description string

		inStartA time.Time
		inEndA   time.Time
		inStartB time.Time
		inEndB   time.Time

		outOverlaps bool
	}{
		{
			description: "should return oberlaps",
			inStartA:    time.Date(2000, 1, 1, 10, 0, 0, 0, time.UTC),
			inEndA:      time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC),
			inStartB:    time.Date(2000, 1, 1, 10, 30, 0, 0, time.UTC),
			inEndB:      time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC),
			outOverlaps: true,
		},
		{
			description: "should not overlaps when its sharing common bounds",
			inStartA:    time.Date(2000, 1, 1, 10, 0, 0, 0, time.UTC),
			inEndA:      time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC),
			inStartB:    time.Date(2000, 1, 1, 11, 0, 0, 0, time.UTC),
			inEndB:      time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			outOverlaps: false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			overlaps := booking.TimeRangesOverlaps(test.inStartA, test.inEndA, test.inStartB, test.inEndB)
			assert.Equal(t, test.outOverlaps, overlaps)
		})
	}
}

func TestLimitTimeRange(t *testing.T) {
	var tests = []struct {
		description string

		lowerLimit time.Time
		upperLimit time.Time
		rangeStart time.Time
		rangeEnd   time.Time

		outLowerRange time.Time
		outUpperRange time.Time
	}{
		{
			description: "should return default time values because range is below of the limits",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 9, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 9, 30, 0, 0, time.UTC),

			outLowerRange: time.Time{},
			outUpperRange: time.Time{},
		},
		{
			description: "should return default time values because range is above the range bounds",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 13, 30, 0, 0, time.UTC),

			outLowerRange: time.Time{},
			outUpperRange: time.Time{},
		},
		{
			description: "should reduce lower range to lower limit",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 8, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),

			outLowerRange: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			outUpperRange: time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),
		},
		{
			description: "should reduce lower range to lower limit",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 8, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),

			outLowerRange: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			outUpperRange: time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),
		},
		{
			description: "should reduce upper range to upper limit",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 19, 30, 0, 0, time.UTC),

			outLowerRange: time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
			outUpperRange: time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
		},
		{
			description: "should return timerange limited to both upper and lower",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 9, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 16, 0, 0, 0, time.UTC),

			outLowerRange: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			outUpperRange: time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
		},
		{
			description: "should return untouched time range if same bounds",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),

			outLowerRange: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			outUpperRange: time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			lowerRange, upperRange := booking.LimitTimeRange(test.lowerLimit, test.upperLimit, test.rangeStart, test.rangeEnd)

			assert.Equal(t, test.outLowerRange, lowerRange, "expected: %s, got: %s", test.outLowerRange, lowerRange)
			assert.Equal(t, test.outUpperRange, upperRange, "expected: %s, got: %s", test.outLowerRange, lowerRange)
		})
	}
}
