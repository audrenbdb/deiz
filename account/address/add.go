package address

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	addressCreater interface {
		CreateClinicianOfficeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
)

type AddAddressUsecase struct {
	AddressCreater addressCreater
}

func (u *AddAddressUsecase) AddClinicianOfficeAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	if address.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	return u.AddressCreater.CreateClinicianOfficeAddress(ctx, address, clinicianID)
}
