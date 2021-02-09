package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianRegistrationCompleteVerifier interface {
		IsClinicianRegistrationComplete(ctx context.Context, email string) (bool, error)
	}
	ClinicianRegistrationCompleter interface {
		CompleteClinicianRegistration(ctx context.Context, email, password string, clinicianID int) error
	}
)

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
	return u.RegistrationCompleter.CompleteClinicianRegistration(ctx, email, password, clinician.ID)
}
