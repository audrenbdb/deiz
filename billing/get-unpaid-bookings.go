package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	unpaidBookingsGetter interface {
		GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error)
	}
)

type GetUnpaidBookingsUsecase struct {
	Getter unpaidBookingsGetter
}

func (u *GetUnpaidBookingsUsecase) GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	return u.Getter.GetUnpaidBookings(ctx, clinicianID)
}
