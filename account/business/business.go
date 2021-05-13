package business

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	updater interface {
		UpdateClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error
	}
	creater interface {
		CreateBusinessAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
	getter interface {
		GetClinicianBusiness(ctx context.Context, clinicianID int) (deiz.Business, error)
	}
	addressUpdater interface {
		UpdateAddress(ctx context.Context, a *deiz.Address) error
	}
)

type Usecase struct {
	BusinessUpdater updater
	AddressCreater  creater
	Getter          getter
	AddressUpdater  addressUpdater
}

func (u *Usecase) EditClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error {
	return u.BusinessUpdater.UpdateClinicianBusiness(ctx, b, clinicianID)
}

func (u *Usecase) SetClinicianBusinessAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	return u.AddressCreater.CreateBusinessAddress(ctx, a, clinicianID)
}

func (u *Usecase) UpdateClinicianBusinessAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	b, err := u.Getter.GetClinicianBusiness(ctx, clinicianID)
	if err != nil {
		return err
	}
	if a.ID != b.Address.ID || a.ID == 0 {
		return deiz.ErrorStructValidation
	}
	return u.AddressUpdater.UpdateAddress(ctx, a)
}
