package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	accountCreater interface {
		CreateClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error
	}
)

type AddAccountUsecase struct {
	AccountCreater accountCreater
}

func (u *AddAccountUsecase) AddAccount(ctx context.Context, account *deiz.ClinicianAccount) error {
	return u.AccountCreater.CreateClinicianAccount(ctx, account)
}
