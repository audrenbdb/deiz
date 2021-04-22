package business

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	updater interface {
		UpdateClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error
	}
)

type UpdateUsecase struct {
	BusinessUpdater updater
}

func (u *UpdateUsecase) EditClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error {
	return u.BusinessUpdater.UpdateClinicianBusiness(ctx, b, clinicianID)
}
