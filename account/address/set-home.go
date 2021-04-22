package address

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	homeAddressSetter interface {
		CreateClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
		SetClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
)

type SetHomeUsecase struct {
	HomeAddressSetter homeAddressSetter
}

func (u *SetHomeUsecase) SetHomeAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	if address.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	if address.IsNotSet() {
		return u.HomeAddressSetter.CreateClinicianHomeAddress(ctx, address, clinicianID)
	}
	return u.HomeAddressSetter.SetClinicianHomeAddress(ctx, address, clinicianID)
}
