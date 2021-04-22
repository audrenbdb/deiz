package clinician

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/valid"
	"strings"
)

type (
	phoneUpdater interface {
		UpdateClinicianPhone(ctx context.Context, phone string, clinicianID int) error
	}
	emailUpdater interface {
		UpdateClinicianEmail(ctx context.Context, phone string, clinicianID int) error
	}
	adeliUpdater interface {
		UpdateClinicianAdeli(ctx context.Context, identifier string, clinicianID int) error
	}
	professionUpdater interface {
		UpdateClinicianProfession(ctx context.Context, profession string, clinicianID int) error
	}
)

type EditUsecase struct {
	PhoneUpdater      phoneUpdater
	EmailUpdater      emailUpdater
	AdeliUpdater      adeliUpdater
	ProfessionUpdater professionUpdater
}

func (u *EditUsecase) EditClinicianPhone(ctx context.Context, phone string, clinicianID int) error {
	if !valid.Phone(phone) {
		return deiz.ErrorStructValidation
	}
	return u.PhoneUpdater.UpdateClinicianPhone(ctx, strings.TrimSpace(phone), clinicianID)
}

func (u *EditUsecase) EditClinicianEmail(ctx context.Context, email string, clinicianID int) error {
	if !valid.Email(email) {
		return deiz.ErrorStructValidation
	}
	return u.EmailUpdater.UpdateClinicianEmail(ctx, strings.ToLower(strings.TrimSpace(email)), clinicianID)
}

func (u *EditUsecase) EditClinicianAdeli(ctx context.Context, adeli string, clinicianID int) error {
	if !adeliValid(adeli) {
		return deiz.ErrorStructValidation
	}
	return u.AdeliUpdater.UpdateClinicianAdeli(ctx, strings.TrimSpace(adeli), clinicianID)
}

func (u *EditUsecase) EditClinicianProfession(ctx context.Context, profession string, clinicianID int) error {
	if !professionNameValid(profession) {
		return deiz.ErrorStructValidation
	}
	return u.ProfessionUpdater.UpdateClinicianProfession(ctx, profession, clinicianID)
}

func adeliValid(adeli string) bool {
	return len(adeli) <= 50 && len(adeli) >= 8
}

func professionNameValid(professionName string) bool {
	return len(professionName) >= 2 && len(professionName) <= 50
}
