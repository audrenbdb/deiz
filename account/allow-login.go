package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	clinicianGetter interface {
		GetClinicianByEmail(ctx context.Context, email string) (deiz.Clinician, error)
	}
	authenticationEnabledChecker interface {
		IsClinicianAuthenticationEnabled(ctx context.Context, email string) (bool, error)
	}
	authenticationEnabler interface {
		EnableClinicianAuthentication(ctx context.Context, clinician *deiz.Clinician, password string) error
	}
)

type AllowLoginUsecase struct {
	ClinicianGetter clinicianGetter
	AuthChecker     authenticationEnabledChecker
	AuthEnabler     authenticationEnabler
}

func (u *AllowLoginUsecase) AllowLogin(ctx context.Context, credentials deiz.Credentials) error {
	if credentials.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	clinician, err := u.ClinicianGetter.GetClinicianByEmail(ctx, credentials.Email)
	if err != nil {
		return err
	}
	enabled, err := u.AuthChecker.IsClinicianAuthenticationEnabled(ctx, credentials.Email)
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}
	return u.AuthEnabler.EnableClinicianAuthentication(ctx, &clinician, credentials.Password)
}
