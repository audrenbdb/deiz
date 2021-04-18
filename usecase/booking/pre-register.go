package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type preRegister struct {
	bookingGetter  clinicianBookingsInTimeRangeGetter
	bookingCreater bookingCreater
}

func NewPreRegisterUsecase(getter clinicianBookingsInTimeRangeGetter,
	creater bookingCreater) *preRegister {
	return &preRegister{
		bookingGetter:  getter,
		bookingCreater: creater,
	}
}

func (r *preRegister) PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error {
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

func (r *preRegister) preRegistrationInvalid(b *deiz.Booking, clinicianID int) bool {
	return b.Confirmed || b.End.Before(b.Start) || b.Blocked || b.Clinician.ID != clinicianID
}
