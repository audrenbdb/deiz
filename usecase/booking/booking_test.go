package booking_test

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase/booking"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockBookingsInTimeRangeGetter struct {
	bookings []deiz.Booking
	err      error
}

type mockBookingCreater struct {
	err error
}

type mockBookingDeleter struct {
	err error
}

func (m *mockBookingsInTimeRangeGetter) GetBookingsInTimeRange(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.Booking, error) {
	return m.bookings, m.err
}

func (m *mockBookingCreater) CreateBooking(ctx context.Context, b *deiz.Booking) error {
	return m.err
}

func (m *mockBookingDeleter) DeleteBooking(ctx context.Context, bookingID, clinicianID int) error {
	return m.err
}

func TestIsBookingValid(t *testing.T) {
	var tests = []struct {
		description string

		inBooking *deiz.Booking
		outValid  bool
	}{
		{
			description: "should be invalid",
			inBooking:   &deiz.Booking{},
		},
		{
			description: "should succeed",
			inBooking: &deiz.Booking{
				Start: time.Now(),
				End:   time.Now().Add(200),
				Clinician: deiz.Clinician{
					ID: 1,
				},
			},
			outValid: true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			valid := booking.IsBookingValid(test.inBooking)
			assert.Equal(t, test.outValid, valid)
		})
	}
}

func TestGetBookingSlots(t *testing.T) {
	var tests = []struct {
		description string

		officeHoursGetter *mockOfficeHoursGetter
		bookingsGetter    *mockBookingsInTimeRangeGetter

		inStart                 time.Time
		inTzName                string
		inDefaultMotiveID       int
		inDefaultMotiveDuration int
		inClinicianID           int

		outBookings []deiz.Booking
		outError    error
	}{
		{
			description:       "should fail to parse timezone",
			officeHoursGetter: &mockOfficeHoursGetter{},
			bookingsGetter:    &mockBookingsInTimeRangeGetter{},
			outError:          deiz.ErrorParsingTimezone,
		},
		{
			description:       "should fail to get clinician office hours",
			bookingsGetter:    &mockBookingsInTimeRangeGetter{},
			officeHoursGetter: &mockOfficeHoursGetter{err: errors.New("failed to get office hours")},

			inTzName: "Europe/Paris",
			outError: errors.New("failed to get office hours"),
		},
		{
			description:       "should fail to get booked slots",
			bookingsGetter:    &mockBookingsInTimeRangeGetter{err: errors.New("failed to retrieve bookings")},
			officeHoursGetter: &mockOfficeHoursGetter{},

			inTzName: "Europe/Paris",
			outError: errors.New("failed to retrieve bookings"),
		},
		{
			description: "should succeed in a booking and a free slot",

			officeHoursGetter: &mockOfficeHoursGetter{
				hours: []deiz.OfficeHours{{StartMn: 0, EndMn: 60, WeekDay: 3}},
			},
			bookingsGetter: &mockBookingsInTimeRangeGetter{
				bookings: []deiz.Booking{
					{
						Start:  time.Date(2021, 2, 9, 23, 45, 0, 0, time.UTC),
						End:    time.Date(2021, 2, 10, 0, 15, 0, 0, time.UTC),
						Motive: deiz.BookingMotive{Duration: 30},
					},
				},
			},

			inStart:                 time.Date(2021, 2, 8, 0, 0, 0, 0, time.UTC),
			inTzName:                "Europe/Paris",
			inDefaultMotiveDuration: 30,

			outBookings: []deiz.Booking{
				{
					Start:  time.Date(2021, 2, 9, 23, 45, 0, 0, time.UTC),
					End:    time.Date(2021, 2, 10, 0, 15, 0, 0, time.UTC),
					Motive: deiz.BookingMotive{Duration: 30},
				},
				{
					Start:  time.Date(2021, 2, 9, 23, 0, 0, 0, time.UTC),
					End:    time.Date(2021, 2, 9, 23, 30, 0, 0, time.UTC),
					Motive: deiz.BookingMotive{Duration: 30},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := booking.Usecase{
				BookingsInTimeRangeGetter: test.bookingsGetter,
				OfficeHoursGetter:         test.officeHoursGetter,
			}
			bookings, err := u.GetBookingSlots(context.Background(), test.inStart, test.inTzName, test.inDefaultMotiveID, test.inDefaultMotiveDuration, test.inClinicianID)
			assert.Equal(t, test.outError, err)
			assert.ElementsMatch(t, test.outBookings, bookings)
		})
	}
}

func TestRemoveOverlappingFreeBookingSlots(t *testing.T) {
	var tests = []struct {
		description string

		inFreeSlots   []deiz.Booking
		inBookedSlots []deiz.Booking
		inSlotsToKeep []deiz.Booking

		outFreeSlots []deiz.Booking
	}{
		{
			description: "should keep only first free slots because 2nd overlaps with booking",
			inFreeSlots: []deiz.Booking{
				{
					Start: time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2021, 2, 12, 0, 30, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2021, 2, 12, 0, 30, 0, 0, time.UTC),
					End:   time.Date(2021, 2, 12, 1, 0, 0, 0, time.UTC),
				},
			},
			inBookedSlots: []deiz.Booking{
				{
					Start: time.Date(2021, 2, 12, 0, 45, 0, 0, time.UTC),
					End:   time.Date(2021, 2, 12, 1, 15, 0, 0, time.UTC),
				},
			},
			inSlotsToKeep: []deiz.Booking{},
			outFreeSlots: []deiz.Booking{
				{
					Start: time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2021, 2, 12, 0, 30, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			freeSlots := booking.RemoveOverlappingFreeBookingSlots(test.inFreeSlots, test.inBookedSlots, test.inSlotsToKeep)
			assert.ElementsMatch(t, test.outFreeSlots, freeSlots)
		})
	}
}

func TestBlockBookingSlot(t *testing.T) {
	var tests = []struct {
		description string

		creater *mockBookingCreater

		inBooking     *deiz.Booking
		inClinicianID int

		outError error
	}{
		{
			description: "should fail to verify booking emptiness",
			inBooking:   &deiz.Booking{Patient: deiz.Patient{ID: 1}},
			outError:    errors.New("booking is not empty"),
		},
		{
			description:   "should not authorize blocking because its not the same user",
			inClinicianID: 1,
			inBooking:     &deiz.Booking{},
			outError:      deiz.ErrorUnauthorized,
		},
		{
			description: "should fail to block booking",
			creater:     &mockBookingCreater{err: errors.New("failed to block")},
			inBooking:   &deiz.Booking{},
			outError:    errors.New("failed to block"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := booking.Usecase{
				Creater: test.creater,
			}
			err := u.BlockBookingSlot(context.Background(), test.inBooking, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}

func TestUnlockBookingSlot(t *testing.T) {
	var tests = []struct {
		description string

		deleter *mockBookingDeleter

		inBookingID    int
		intClinicianID int

		outError error
	}{
		{
			description: "should fail to delete",
			deleter:     &mockBookingDeleter{err: errors.New("fail to delete")},
			outError:    errors.New("fail to delete"),
		},
		{
			description: "should succeed",
			deleter:     &mockBookingDeleter{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := booking.Usecase{
				Deleter: test.deleter,
			}
			err := u.UnlockBookingSlot(context.Background(), test.inBookingID, test.intClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}
