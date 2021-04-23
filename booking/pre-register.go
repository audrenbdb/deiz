package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type PreRegisterUsecase struct {
	BookingGetter  clinicianBookingsInTimeRangeGetter
	BookingCreater bookingCreater
}

//PreRegisterBooking locks a given slot to be completed later by adding a patient or settings different details.
//Its similar to registration but booking status wont be confirmed and mail reminder wont be send
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
