package booking

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type repo struct {
	RecurrentBookingsGetter
}

type RecurrentBookingsGetter interface {
	GetRecurrentBookingsByClinicianID(ctx context.Context, clinicianID int) ([]deiz.Booking, error)
}

type psql struct{}

func (r *psql) GetRecurrentBookingsByClinicianID(ctx context.Context, clinicianID int) ([]deiz.Booking, error) {
	return nil, nil
}
