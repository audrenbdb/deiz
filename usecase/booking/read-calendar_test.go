package booking

import (
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstraintTimeRangeWithinLimit(t *testing.T) {
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
			description: "should reduce lower range to lower limit",
			lowerLimit:  time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			upperLimit:  time.Date(2010, 1, 1, 15, 0, 0, 0, time.UTC),
			rangeStart:  time.Date(2010, 1, 1, 8, 0, 0, 0, time.UTC),
			rangeEnd:    time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),

			outLowerRange: time.Date(2010, 1, 1, 10, 0, 0, 0, time.UTC),
			outUpperRange: time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			timeRange := constraintTimeRangeWithinLimit(test.lowerLimit, test.upperLimit, timeRange{test.rangeStart, test.rangeEnd})

			assert.Equal(t, test.outLowerRange, timeRange.start, "expected: %s, got: %s", test.outLowerRange, timeRange.start)
			assert.Equal(t, test.outUpperRange, timeRange.end, "expected: %s, got: %s", test.outLowerRange, timeRange.end)
		})
	}
}

func TestConvertOfficeHoursToTimeRange(t *testing.T) {
	u := NewCalendarReaderUsecase(CalendarReaderDeps{Loc: time.UTC})

	var tests = []struct {
		description string

		start time.Time
		end   time.Time
		h     deiz.OfficeHours

		outputTimerange [2]time.Time
	}{
		{
			description: "should return time range limited by start value",

			start: time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 1, 8, 10, 0, 0, 0, time.UTC),
			h:     deiz.OfficeHours{StartMn: 540, EndMn: 720, WeekDay: 5},

			outputTimerange: [2]time.Time{
				time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			description: "should return time range limited by end value",

			start: time.Date(2021, 1, 2, 10, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 1, 8, 10, 0, 0, 0, time.UTC),
			h:     deiz.OfficeHours{StartMn: 480, EndMn: 720, WeekDay: 5},

			outputTimerange: [2]time.Time{
				time.Date(2021, 1, 8, 8, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 8, 10, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			timeRange := u.convertOfficeHoursToTimeRange(timeRange{test.start, test.end}, test.h)
			assert.Equal(t, test.outputTimerange, timeRange, "expected : %s, got : %s", test.outputTimerange, timeRange)
		})
	}

}
