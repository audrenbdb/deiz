package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type PreRegisterUsecase struct {
	BookingGetter  clinicianBookingsInTimeRangeGetter
	BookingCreater bookingCreater
}

func (r *PreRegisterUsecase) PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	if r.preRegistrationInvalid(b, clinicianID) {
		return deiz.ErrorStructValidation
	}
	available, err := bookingSlotAvailable(ctx, b, r.BookingGetter)
	if err != nil {
		return err
	}
	if !available {
		return deiz.ErrorBookingSlotAlreadyFilled
	}
	return r.BookingCreater.CreateBooking(ctx, b)
}

func (r *PreRegisterUsecase) preRegistrationInvalid(b *deiz.Booking, clinicianID int) bool {
	return b.Confirmed || b.End.Before(b.Start) || b.Blocked || b.Clinician.ID != clinicianID
}
