package booking_test

import (
	"github.com/audrenbdb/deiz"
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

func TestGetTimeRangesFromBookings(t *testing.T) {
	var tests = []struct {
		description string
		inBookings  []deiz.Booking
		outRanges   [][2]time.Time
	}{
		{
			description: "should return a list of two time ranges",
			inBookings: []deiz.Booking{
				{
					Start: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
					End:   time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
					End:   time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
				},
			},
			outRanges: [][2]time.Time{
				{
					time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
					time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
				},
				{
					time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
					time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			timeRanges := booking.GetTimeRangesFromBookings(test.inBookings, [][2]time.Time{})
			assert.ElementsMatch(t, test.outRanges, timeRanges)
		})
	}
}

func TestSortBookingsByStart(t *testing.T) {
	t.Run("should sort two bookings", func(t *testing.T) {
		inBookings := []deiz.Booking{
			{
				Start: time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
			},
			{
				Start: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
				End:   time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
			},
		}
		outBookings := []deiz.Booking{
			{
				Start: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
				End:   time.Date(2010, 1, 1, 11, 0, 0, 0, time.UTC),
			}, {
				Start: time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
			},
		}
		bookings := booking.SortBookingsByStart(inBookings)
		assert.Equal(t, outBookings, bookings)
	})
}

func TestGetTimeRangesNotOverlapping(t *testing.T) {
	var tests = []struct {
		description string

		inDuration           int
		inLowerLimit         time.Time
		inUpperLimit         time.Time
		inRangesToNotOverlap [][2]time.Time

		outRanges [][2]time.Time
	}{
		{
			description:  "should return two 30mn time ranges that would fit lower and upper limit",
			inDuration:   30,
			inLowerLimit: time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
			inUpperLimit: time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
			outRanges: [][2]time.Time{
				{
					time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
					time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),
				},
				{
					time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),
					time.Date(2010, 1, 1, 13, 00, 0, 0, time.UTC),
				},
			},
		},
		{
			description:  "should return one 30mn time range before a time range blocked",
			inDuration:   30,
			inLowerLimit: time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
			inUpperLimit: time.Date(2010, 1, 1, 13, 0, 0, 0, time.UTC),
			inRangesToNotOverlap: [][2]time.Time{
				{
					time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),
					time.Date(2010, 1, 1, 13, 00, 0, 0, time.UTC),
				},
			},
			outRanges: [][2]time.Time{
				{
					time.Date(2010, 1, 1, 12, 0, 0, 0, time.UTC),
					time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ranges := booking.GetTimeRangesNotOverLapping(test.inDuration, test.inLowerLimit, test.inUpperLimit, test.inRangesToNotOverlap, [][2]time.Time{})
			assert.ElementsMatch(t, test.outRanges, ranges)
		})
	}
}
