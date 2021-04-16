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

func phoneValid(phone string) bool {
	return len(phone) >= 10 && len(phone) <= 20
}

func emailValid(email string) bool {
	return len(email) >= 3 && len(email) <= 254 && emailRegex.MatchString(email)
}

func (u *Usecase) EditClinicianPhone(ctx context.Context, phone string, clinicianID int) error {
	if !phoneValid(phone) {
		return deiz.ErrorStructValidation
	}
	return u.ClinicianPhoneUpdater.UpdateClinicianPhone(ctx, phone, clinicianID)
}

func (u *Usecase) EditClinicianEmail(ctx context.Context, email string, clinicianID int) error {
	if !emailValid(email) {
		return deiz.ErrorStructValidation
	}
	return u.ClinicianEmailUpdater.UpdateClinicianEmail(ctx, email, clinicianID)
}

func adeliValid(adeli string) bool {
	return len(adeli) <= 50 && len(adeli) >= 8
}

func professionNameValid(professionName string) bool {
	return len(professionName) >= 2 && len(professionName) <= 50
}

func (u *Usecase) EditClinicianAdeli(ctx context.Context, identifier string, clinicianID int) error {
	if !adeliValid(identifier) {
		return deiz.ErrorUnauthorized
	}
	return u.ClinicianAdeliUpdater.UpdateClinicianAdeli(ctx, identifier, clinicianID)
}

func (u *Usecase) EditClinicianProfession(ctx context.Context, profession string, clinicianID int) error {
	if !professionNameValid(profession) {
		return deiz.ErrorUnauthorized
	}
	return u.ClinicianProfessionUpdater.UpdateClinicianProfession(ctx, profession, clinicianID)
}
