package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	bookingCreater interface {
		CreateBooking(ctx context.Context, b *deiz.Booking) error
	}
	bookingMailer interface {
		MailBookingToClinician(b *deiz.Booking) error
		MailBookingToPatient(b *deiz.Booking) error
	}
)

func bookingsOverlap(booking1, booking2 *deiz.Booking) bool {
	return booking1.Start.Before(booking2.End) && booking2.Start.Before(booking1.End)
}

func bookingSlotAvailable(ctx context.Context, b *deiz.Booking, getter bookingGetter, loc *time.Location) (bool, error) {
	overlapExistingBookings, err := bookingOverlapExistingBookings(ctx, b, getter)
	if err != nil {
		return false, err
	}
	if overlapExistingBookings {
		return false, nil
	}
	overlapRecurrentBookings, err := bookingOverlapRecurrentBookings(ctx, b, getter, loc)
	if err != nil {
		return false, err
	}
	if overlapRecurrentBookings {
		return false, nil
	}
	return true, nil
}

func bookingOverlapRecurrentBookings(ctx context.Context, b *deiz.Booking, getter bookingGetter, loc *time.Location) (bool, error) {
	recurrentBookings, err := getter.GetClinicianWeeklyRecurrentBookings(ctx, b.Clinician.ID)
	if err != nil {
		return true, err
	}
	if recurrentBookingsOverlapTimeRange(timeRange{start: b.Start, end: b.End}, recurrentBookings, loc) {
		return true, nil
	}
	return false, nil
}

func bookingOverlapExistingBookings(ctx context.Context, b *deiz.Booking, getter bookingGetter) (bool, error) {
	bookings, err := getter.GetNonRecurrentClinicianBookingsInTimeRange(ctx, b.Start, b.End, b.Clinician.ID)
	if err != nil {
		return false, err
	}
	for _, booking := range bookings {
		if bookingsOverlap(b, &booking) && booking.ID != b.ID {
			return true, nil
		}
	}
	return false, nil
}

func recurrentBookingsOverlapTimeRange(tr timeRange, rb []deiz.Booking, loc *time.Location) bool {
	for _, r := range rb {
		if recurrentBookingInTimeRange(tr, r, loc) {
			return true
		}
	}
	return false
}

func recurrentBookingInTimeRange(tr timeRange, r deiz.Booking, loc *time.Location) bool {
	tr = convertCalEventToTimeRange(tr, calEvent{
		weekday: int(r.Start.In(loc).Weekday()),
		startMn: r.Start.In(loc).Hour()*60 + r.Start.In(loc).Minute(),
		endMn:   r.End.In(loc).Hour()*60 + r.End.In(loc).Minute(),
		loc:     loc,
	}, true)
	return !tr.isNull()
}

func filterNonRecurrentBookings(bookings []deiz.Booking) []deiz.Booking {
	nonRecurrentBookings := []deiz.Booking{}
	for _, b := range bookings {
		if b.Recurrence == deiz.NoRecurrence {
			nonRecurrentBookings = append(nonRecurrentBookings, b)
		}
	}
	return nonRecurrentBookings
}

func filterConfirmedBookings(bookings []deiz.Booking) []deiz.Booking {
	confirmedBookings := []deiz.Booking{}
	for _, b := range bookings {
		if b.Confirmed && b.BookingType == deiz.AppointmentBooking {
			confirmedBookings = append(confirmedBookings, b)
		}
	}
	return confirmedBookings
}

func areBookingsValid(bookings []*deiz.Booking, clinicianID int) bool {
	for _, b := range bookings {
		if b.IsInvalid(clinicianID) {
			return false
		}
	}
	return true
}

func areBookingsInvalid(bookings []*deiz.Booking, clinicianID int) bool {
	return !areBookingsValid(bookings, clinicianID)
}
