package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type PreRegisterUsecase struct {
	BookingGetter  bookingGetter
	BookingCreater bookingCreater
	Loc            *time.Location
}

//PreRegisterBooking locks a given slot to be completed later by adding a patient or settings different details.
//Its similar to registration but booking status wont be confirmed and mail reminder wont be send
func (r *PreRegisterUsecase) PreRegisterBookings(ctx context.Context, bookings []*deiz.Booking, clinicianID int) error {
	return registerBookings(ctx,
		registrationDependencies{creater: r.BookingCreater, getter: r.BookingGetter, loc: r.Loc},
		bookings, clinicianID, false, false)
}
