package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type PreRegister struct {
	bookingGetter  clinicianBookingsInTimeRangeGetter
	bookingCreater bookingCreater
}

type PreRegisterDeps struct {
	BookingGetter  clinicianBookingsInTimeRangeGetter
	BookingCreater bookingCreater
}

func NewPreRegisterUsecase(deps PreRegisterDeps) *PreRegister {
	return &PreRegister{
		bookingGetter:  deps.BookingGetter,
		bookingCreater: deps.BookingCreater,
	}
}

func (r *PreRegister) PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error {
	if r.preRegistrationInvalid(b, clinicianID) {
		return deiz.ErrorStructValidation
	}
	available, err := bookingSlotAvailable(ctx, b, r.bookingGetter)
	if err != nil {
		return err
	}
	if !available {
		return deiz.ErrorBookingSlotAlreadyFilled
	}
	return r.bookingCreater.CreateBooking(ctx, b)
}

func (r *PreRegister) preRegistrationInvalid(b *deiz.Booking, clinicianID int) bool {
	return b.Confirmed || b.End.Before(b.Start) || b.Blocked || b.Clinician.ID != clinicianID
}
