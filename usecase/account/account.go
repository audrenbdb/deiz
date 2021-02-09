package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianAccountCreater interface {
		CreateClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error
	}
	ClinicianAccountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error)
	}
)

func (u *Usecase) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return u.AccountGetter.GetClinicianAccount(ctx, clinicianID)
}

func (u *Usecase) AddClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error {
	return u.AccountCreater.CreateClinicianAccount(ctx, account)
}
