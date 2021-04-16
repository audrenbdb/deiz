package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	OfficeHoursCreater interface {
		CreateOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error
	}
	OfficeHoursDeleter interface {
		DeleteOfficeHours(ctx context.Context, hoursID, clinicianID int) error
	}
)

func (u *Usecase) AddOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error {
	if h.IsInvalid() {
		return deiz.ErrorUnauthorized
	}
	if h.Address.IsSet() {
		owns, err := u.AddressOwnerShipVerifier.IsAddressToClinician(ctx, &deiz.Address{ID: h.Address.ID}, clinicianID)
		if err != nil {
			return err
		}
		if !owns {
			return deiz.ErrorUnauthorized
		}
	}
	return u.OfficeHoursCreater.CreateOfficeHours(ctx, h, clinicianID)
}

func (u *Usecase) RemoveOfficeHours(ctx context.Context, hoursID, clinicianID int) error {
	return u.OfficeHoursDeleter.DeleteOfficeHours(ctx, hoursID, clinicianID)
}
