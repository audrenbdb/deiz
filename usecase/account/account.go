package account

import (
	"context"
	"errors"
	"fmt"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianAccountCreater interface {
		CreateClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error
	}
	ClinicianAccountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error)
	}
	PublicDataGetter interface {
		GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccountPublicData, error)
	}
	ClinicianRegistrationCompleteVerifier interface {
		IsClinicianRegistrationComplete(ctx context.Context, email string) (bool, error)
	}
	ClinicianRegistrationCompleter interface {
		CompleteClinicianRegistration(ctx context.Context, c *deiz.Clinician, password string, clinicianID int) error
	}
)

func (u *Usecase) GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccountPublicData, error) {
	return u.PublicDataGetter.GetClinicianAccountPublicData(ctx, clinicianID)
}

func (u *Usecase) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	account, err := u.AccountGetter.GetClinicianAccount(ctx, clinicianID)
	if err != nil {
		return deiz.ClinicianAccount{}, err
	}
	return account, nil
}

func (u *Usecase) AddClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error {
	return u.AccountCreater.CreateClinicianAccount(ctx, account)
}

func (u *Usecase) EnsureClinicianRegistrationComplete(ctx context.Context, email, password string) error {
	if !credentialsValid(email, password) {
		return deiz.ErrorStructValidation
	}
	clinician, err := u.ClinicianGetterByEmail.GetClinicianByEmail(ctx, email)
	if err != nil {
		return errors.New("ce clinician n'existe pas")
	}
	isComplete, err := u.RegistrationVerifier.IsClinicianRegistrationComplete(ctx, email)
	if err != nil {
		return fmt.Errorf("unable to check is registration is complete: %s", err)
	}
	if !isComplete {
		return u.RegistrationCompleter.CompleteClinicianRegistration(ctx, &clinician, password, clinician.ID)
	}
	return nil
}

func credentialsValid(email, password string) bool {
	return len(email) >= 5 && len(password) >= 6
}
