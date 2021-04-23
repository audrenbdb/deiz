package address

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	addressDeleter interface {
		DeleteAddress(ctx context.Context, addressID int) error
	}
)

type DeleteAddressUsecase struct {
	AccountGetter  accountGetter
	AddressDeleter addressDeleter
}

func (u *DeleteAddressUsecase) DeleteAddress(ctx context.Context, addressID int, cred deiz.Credentials) error {
	authorized, err := isAddressToClinician(ctx, addressID, cred.UserID, u.AccountGetter)
	if err != nil {
		return err
	}
	if !authorized {
		return deiz.ErrorUnauthorized
	}
	return u.AddressDeleter.DeleteAddress(ctx, addressID)
}
