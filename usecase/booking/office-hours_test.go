package booking_test

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase/booking"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockOfficeHoursGetter struct {
	hours []deiz.OfficeHours
	err   error
}

func (m *mockOfficeHoursGetter) GetClinicianOfficeHours(ctx context.Context, clinicianID int) ([]deiz.OfficeHours, error) {
	return m.hours, m.err
}

func TestGetOfficeHoursTimeRange(t *testing.T) {
	paris, _ := time.LoadLocation("Europe/Paris")
	utc := time.UTC
	officeHours := deiz.OfficeHours{
		StartMn: 0,
		EndMn:   120,
		WeekDay: 3,
	}

	var tests = []struct {
		description string

		inAnchor time.Time
		inEnd    time.Time
		inHours  deiz.OfficeHours
		inLoc    *time.Location

		outLowerRange time.Time
		outUpperRange time.Time
	}{
		{
			description: "should return nil time values because office hours week day is not in time range",

			inAnchor: time.Date(2021, 2, 12, 0, 0, 0, 0, utc),
			inEnd:    time.Date(2021, 2, 14, 0, 0, 0, 0, utc),
			inHours:  officeHours,
			inLoc:    utc,

			outLowerRange: time.Time{},
			outUpperRange: time.Time{},
		},
		{
			description: "should return time range that matches office hours in paris location",

			inAnchor: time.Date(2021, 2, 8, 0, 0, 0, 0, utc),
			inEnd:    time.Date(2021, 2, 11, 0, 0, 0, 0, utc),
			inHours:  officeHours,
			inLoc:    paris,

			outLowerRange: time.Date(2021, 2, 9, 23, 0, 0, 0, utc),
			outUpperRange: time.Date(2021, 2, 10, 1, 0, 0, 0, utc),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			lowerRange, upperRange := booking.GetOfficeHoursTimeRange(test.inAnchor, test.inEnd, test.inHours, test.inLoc)
			assert.Equal(t, test.outLowerRange, lowerRange, "got: %s, expected: %s", lowerRange, test.outLowerRange)
			assert.Equal(t, test.outUpperRange, upperRange, "got: %s, expected: %s", upperRange, test.outUpperRange)
		})
	}
}
