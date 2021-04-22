package patient

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	AddressCreater interface {
		CreatePatientAddress(ctx context.Context, a *deiz.Address, patientID int) error
	}
	AddressUpdater interface {
		UpdateAddress(ctx context.Context, a *deiz.Address) error
	}
	ClinicianBoundChecker interface {
		IsPatientTiedToClinician(ctx context.Context, p *deiz.Patient, clinicianID int) (bool, error)
	}
)

func IsAddressValid(a *deiz.Address) bool {
	if len(a.Line) < 2 || a.PostCode < 10000 || len(a.City) < 2 {
		return false
	}
	return true
}

func (u *Usecase) AddPatientAddress(ctx context.Context, a *deiz.Address, patientID int, clinicianID int) error {
	if !IsAddressValid(a) {
		return deiz.ErrorStructValidation
	}
	bound, err := u.ClinicianBoundChecker.IsPatientTiedToClinician(ctx, &deiz.Patient{ID: patientID}, clinicianID)
	if err != nil {
		return err
	}
	if !bound {
		return deiz.ErrorUnauthorized
	}
	return u.AddressCreater.CreatePatientAddress(ctx, a, patientID)
}

func (u *Usecase) EditPatientAddress(ctx context.Context, a *deiz.Address, patientID int, clinicianID int) error {
	if !IsAddressValid(a) {
		return deiz.ErrorStructValidation
	}
	bound, err := u.ClinicianBoundChecker.IsPatientTiedToClinician(ctx, &deiz.Patient{ID: patientID}, clinicianID)
	if err != nil {
		return err
	}
	if !bound {
		return deiz.ErrorUnauthorized
	}
	return u.AddressUpdater.UpdateAddress(ctx, a)
}
