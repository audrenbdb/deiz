package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianBusinessUpdater interface {
		UpdateClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error
	}
)

func (u *Usecase) EditClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error {
	return u.BusinessUpdater.UpdateClinicianBusiness(ctx, b, clinicianID)
}
