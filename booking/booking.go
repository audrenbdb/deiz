package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	clinicianBookingsInTimeRangeGetter interface {
		GetClinicianBookingsInTimeRange(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.Booking, error)
	}
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

func bookingSlotAvailable(ctx context.Context, b *deiz.Booking, getter clinicianBookingsInTimeRangeGetter) (bool, error) {
	bookings, err := getter.GetClinicianBookingsInTimeRange(ctx, b.Start, b.End, b.Clinician.ID)
	if err != nil {
		return false, err
	}
	for _, booking := range bookings {
		if bookingsOverlap(b, &booking) {
			return false, nil
		}
	}
	return true, nil
}
