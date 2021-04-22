package patient

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	Searcher interface {
		SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error)
	}
	Creater interface {
		CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	Updater interface {
		UpdatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	GetterByEmail interface {
		GetPatientByEmail(ctx context.Context, email string, clinicianID int) (deiz.Patient, error)
	}
)

func (u *Usecase) SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error) {
	return u.Searcher.SearchPatient(ctx, search, clinicianID)
}

func (u *Usecase) AddPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	if p.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	existingPatient, err := u.GetterByEmail.GetPatientByEmail(ctx, p.Email, clinicianID)
	if err != nil {
		return err
	}
	if existingPatient.ID != 0 {
		p.ID = existingPatient.ID
		return nil
	}
	return u.Creater.CreatePatient(ctx, p, clinicianID)
}

func (u *Usecase) EditPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	if p.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	return u.Updater.UpdatePatient(ctx, p, clinicianID)
}
