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

func TestFillOfficeHoursWithFreeBookingSlots(t *testing.T) {
	paris, _ := time.LoadLocation("Europe/Paris")
	utc := time.UTC
	officeHours := []deiz.OfficeHours{
		{
			StartMn: 0,
			EndMn:   60,
			WeekDay: 3,
		},
	}
	defaultMotive := deiz.BookingMotive{Duration: 30}

	var tests = []struct {
		description string

		inStart     time.Time
		inEnd       time.Time
		inClinician deiz.Clinician
		inHours     []deiz.OfficeHours
		inBookings  []deiz.Booking
		inMotive    deiz.BookingMotive
		inLoc       *time.Location

		outBookings []deiz.Booking
	}{
		{
			description: "should return two 30mn booking slots",
			inStart:     time.Date(2021, 2, 8, 0, 0, 0, 0, utc),
			inEnd:       time.Date(2021, 2, 11, 0, 0, 0, 0, utc),
			inHours:     officeHours,
			inBookings:  []deiz.Booking{},
			inMotive:    defaultMotive,
			inLoc:       paris,

			outBookings: []deiz.Booking{
				{
					Start:  time.Date(2021, 2, 9, 23, 0, 0, 0, utc),
					End:    time.Date(2021, 2, 9, 23, 30, 0, 0, utc),
					Motive: defaultMotive,
				},
				{
					Start:  time.Date(2021, 2, 9, 23, 30, 0, 0, utc),
					End:    time.Date(2021, 2, 10, 0, 0, 0, 0, utc),
					Motive: defaultMotive,
				},
			},
		},
		{
			description: "should return nothing because no availability",
			inStart:     time.Date(2021, 2, 8, 0, 0, 0, 0, utc),
			inEnd:       time.Date(2021, 2, 11, 0, 0, 0, 0, utc),
			inHours:     []deiz.OfficeHours{},
			inBookings:  []deiz.Booking{},
			inMotive:    defaultMotive,
			inLoc:       paris,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			bookings := booking.FillOfficeHoursWithFreeBookingSlots(test.inStart, test.inEnd, test.inClinician, test.inHours, test.inBookings, test.inMotive, test.inLoc)
			assert.ElementsMatch(t, test.outBookings, bookings)
		})
	}
}

func TestFillTimeRangeWithFreeBookingSlots(t *testing.T) {
	utc := time.UTC
	defaultMotive := deiz.BookingMotive{Duration: 30}

	var tests = []struct {
		description string

		inStart     time.Time
		inEnd       time.Time
		inClinician deiz.Clinician
		inBookings  []deiz.Booking
		inAddress   deiz.Address
		inMotive    deiz.BookingMotive

		outBooking []deiz.Booking
	}{
		{
			description: "should return two 30mn booking slots",
			inStart:     time.Date(2021, 2, 8, 0, 0, 0, 0, utc),
			inEnd:       time.Date(2021, 2, 8, 1, 0, 0, 0, utc),
			inBookings:  []deiz.Booking{},
			inMotive:    defaultMotive,

			outBooking: []deiz.Booking{
				{
					Start:  time.Date(2021, 2, 8, 0, 0, 0, 0, utc),
					End:    time.Date(2021, 2, 8, 0, 30, 0, 0, utc),
					Motive: defaultMotive,
				},
				{
					Start:  time.Date(2021, 2, 8, 0, 30, 0, 0, utc),
					End:    time.Date(2021, 2, 8, 1, 0, 0, 0, utc),
					Motive: defaultMotive,
				},
			},
		},
		{
			description: "should return nil because there is not enough time for a booking to fit in",
			inStart:     time.Date(2021, 2, 8, 0, 0, 0, 0, utc),
			inEnd:       time.Date(2021, 2, 8, 0, 20, 0, 0, utc),
			inBookings:  []deiz.Booking{},
			inMotive:    defaultMotive,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			bookings := booking.FillTimeRangeWithFreeBookingSlots(test.inStart, test.inEnd, test.inClinician, test.inBookings, test.inAddress, test.inMotive)
			assert.ElementsMatch(t, test.outBooking, bookings)
		})
	}
}
