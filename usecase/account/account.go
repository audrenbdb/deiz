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
	ClinicianRegistrationCompleteVerifier interface {
		IsClinicianRegistrationComplete(ctx context.Context, email string) (bool, error)
	}
	ClinicianRegistrationCompleter interface {
		CompleteClinicianRegistration(ctx context.Context, c *deiz.Clinician, password string, clinicianID int) error
	}
)

func (u *Usecase) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return u.AccountGetter.GetClinicianAccount(ctx, clinicianID)
}

func (u *Usecase) AddClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error {
	return u.AccountCreater.CreateClinicianAccount(ctx, account)
}

//EnsureClinicianRegistration checks if account exists and is registered with email and password
//if not, complete it
func (u *Usecase) EnsureClinicianRegistrationComplete(ctx context.Context, email, password string) error {
	if len(email) < 5 || len(password) < 6 {
		return deiz.ErrorStructValidation
	}
	clinician, err := u.ClinicianGetterByEmail.GetClinicianByEmail(ctx, email)
	if err != nil {
		return err
	}
	isRegistered, err := u.RegistrationVerifier.IsClinicianRegistrationComplete(ctx, email)
	if err != nil {
		return err
	}
	if isRegistered {
		return nil
	}
	return u.RegistrationCompleter.CompleteClinicianRegistration(ctx, &clinician, password, clinician.ID)
}
