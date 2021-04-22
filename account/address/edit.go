package address

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func (u *EditAddressUsecase) EditAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	if address.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	allowedToEdit, err := isAddressToClinician(ctx, address.ID, clinicianID, u.AccountGetter)
	if err != nil {
		return err
	}
	if !allowedToEdit {
		return deiz.ErrorUnauthorized
	}
	return u.AddressUpdater.UpdateAddress(ctx, address)
}

type (
	addressUpdater interface {
		UpdateAddress(ctx context.Context, address *deiz.Address) error
	}
	accountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error)
	}
)

type EditAddressUsecase struct {
	AddressUpdater addressUpdater
	AccountGetter  accountGetter
}
