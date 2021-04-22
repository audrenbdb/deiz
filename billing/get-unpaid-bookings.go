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

type GetUnpaidBookings struct {
	getter unpaidBookingsGetter
}

func NewGetUnpaidBookings(getter unpaidBookingsGetter) *GetUnpaidBookings {
	return &GetUnpaidBookings{
		getter: getter,
	}
}

func (u *GetUnpaidBookings) GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	return u.getter.GetUnpaidBookings(ctx, clinicianID)
}
