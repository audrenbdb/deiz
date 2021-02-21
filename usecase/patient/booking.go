package patient

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	BookingsGetter interface {
		GetPatientBookings(ctx context.Context, clinicianID int, patientID int) ([]deiz.Booking, error)
	}
)

func (u *Usecase) GetPatientBookings(ctx context.Context, clinicianID, patientID int) ([]deiz.Booking, error) {
	return u.BookingsGetter.GetPatientBookings(ctx, clinicianID, patientID)
}
