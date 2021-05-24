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
func (r *PreRegisterUsecase) PreRegisterBookings(ctx context.Context, bookings []*deiz.Booking, clinicianID int) error {
	if r.preRegistrationsInvalid(bookings, clinicianID) {
		return deiz.ErrorStructValidation
	}
	for _, b := range bookings {
		available, err := bookingSlotAvailable(ctx, b, r.BookingGetter)
		if err != nil {
			return err
		}
		if !available {
			return deiz.ErrorBookingSlotAlreadyFilled
		}
		if err := r.BookingCreater.CreateBooking(ctx, b); err != nil {
			return err
		}
	}
	return nil
}

func (r *PreRegisterUsecase) preRegistrationsInvalid(bookings []*deiz.Booking, clinicianID int) bool {
	for _, b := range bookings {
		if b.BookingType == deiz.EventBooking {
			b.ToEvent()
		}
		if b.Confirmed || b.End.Before(b.Start) || b.BookingType == deiz.BlockedBooking || b.Clinician.ID != clinicianID {
			return true
		}
	}
	return false
}
