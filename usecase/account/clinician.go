package account

import (
	"context"
	"github.com/audrenbdb/deiz"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type (
	ClinicianGetterByEmail interface {
		GetClinicianByEmail(ctx context.Context, email string) (deiz.Clinician, error)
	}
	ClinicianPhoneUpdater interface {
		UpdateClinicianPhone(ctx context.Context, phone string, clinicianID int) error
	}
	ClinicianEmailUpdater interface {
		UpdateClinicianEmail(ctx context.Context, email string, clinicianID int) error
	}
	ClinicianAdeliUpdater interface {
		UpdateClinicianAdeli(ctx context.Context, identifier string, clinicianID int) error
	}
	ClinicianProfessionUpdater interface {
		UpdateClinicianProfession(ctx context.Context, profession string, clinicianID int) error
	}
)

func (u *Usecase) EditClinicianPhone(ctx context.Context, phone string, clinicianID int) error {
	if len(phone) < 10 || len(phone) > 20 {
		return deiz.ErrorStructValidation
	}
	return u.ClinicianPhoneUpdater.UpdateClinicianPhone(ctx, phone, clinicianID)
}

func (u *Usecase) EditClinicianEmail(ctx context.Context, email string, clinicianID int) error {
	if (len(email) < 3 && len(email) > 254) || !emailRegex.MatchString(email) {
		return deiz.ErrorStructValidation
	}
	return u.ClinicianEmailUpdater.UpdateClinicianEmail(ctx, email, clinicianID)
}

func (u *Usecase) EditClinicianAdeli(ctx context.Context, identifier string, clinicianID int) error {
	if len(identifier) > 50 || len(identifier) < 8 {
		return deiz.ErrorUnauthorized
	}
	return u.ClinicianAdeliUpdater.UpdateClinicianAdeli(ctx, identifier, clinicianID)
}

func (u *Usecase) EditClinicianProfession(ctx context.Context, profession string, clinicianID int) error {
	if len(profession) < 2 || len(profession) > 50 {
		return deiz.ErrorUnauthorized
	}
	return u.ClinicianProfessionUpdater.UpdateClinicianProfession(ctx, profession, clinicianID)
}
