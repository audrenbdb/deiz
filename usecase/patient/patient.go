package patient

import (
	"context"
	"github.com/audrenbdb/deiz"
	"regexp"
)

type (
	Searcher interface {
		SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error)
	}
	Creater interface {
		CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
)

func IsPatientValid(p *deiz.Patient) bool {
	if len(p.Name) < 2 {
		return false
	}
	if len(p.Surname) < 2 {
		return false
	}
	if len(p.Phone) < 10 {
		return false
	}
	r := regexp.MustCompile("^\\S+@\\S+$")
	if !r.MatchString(p.Email) {
		return false
	}
	return true
}

func (u *Usecase) SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error) {
	return u.Searcher.SearchPatient(ctx, search, clinicianID)
}

func (u *Usecase) AddPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	if !IsPatientValid(p) {
		return deiz.ErrorStructValidation
	}
	return u.Creater.CreatePatient(ctx, p, clinicianID)
}
