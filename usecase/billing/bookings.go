package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	UnpaidBookingsGetter interface {
		GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error)
	}
)

func (u *Usecase) GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	return u.UnpaidBookingsGetter.GetUnpaidBookings(ctx, clinicianID)
}
