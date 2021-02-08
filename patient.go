package deiz

import (
	"context"
	"regexp"
)

//Patient uses the application to book clinician appointment
type Patient struct {
	ID      int     `json:"id" validator:"required"`
	Name    string  `json:"name" validator:"required"`
	Surname string  `json:"surname" validator:"required"`
	Phone   string  `json:"phone" validator:"required,min=10"`
	Email   string  `json:"email" validator:"required,email"`
	Address Address `json:"address" validator:"required"`
}

func (p *Patient) IsValid() bool {
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

type (
	patientEditer interface {
		EditPatient(ctx context.Context, p *Patient, clinicianID int) error
	}
	patientsSearcher interface {
		SearchPatients(ctx context.Context, search string, clinicianID int) ([]Patient, error)
	}
	patientsCounter interface {
		CountPatients(ctx context.Context, clinicianID int) (int, error)
	}
	patientRemover interface {
		RemovePatient(ctx context.Context, p *Patient, clinicianID int) error
	}
	patientCreater interface {
		CreatePatient(ctx context.Context, p *Patient, clinicianID int) error
	}
	patientAddressEditer interface {
		EditPatientAddress(ctx context.Context, p *Patient, clinicianID int) error
	}
)

type (
	//EditPatient edit patient data
	EditPatient func(ctx context.Context, p *Patient, clinicianID int) error
	//SearchPatients search for patients with a given query
	SearchPatients func(ctx context.Context, query string, clinicianID int) ([]Patient, error)
	//CountPatients count total number of patients for a given clinician
	CountPatients func(ctx context.Context, clinicianID int) (int, error)
	//RemovePatient remove a given patient of a given clinician
	RemovePatient func(ctx context.Context, p *Patient, clinicianID int) error
	//EditPatientAddress edit patient address
	EditPatientAddress func(ctx context.Context, p *Patient, clinicianID int) error
)

func editPatientFunc(editer patientEditer) EditPatient {
	return func(ctx context.Context, p *Patient, clinicianID int) error {
		return editer.EditPatient(ctx, p, clinicianID)
	}
}

func searchPatientsFunc(searcher patientsSearcher) SearchPatients {
	return func(ctx context.Context, query string, clinicianID int) ([]Patient, error) {
		return nil, nil
	}
}

func countPatientsFunc(counter patientsCounter) CountPatients {
	return func(ctx context.Context, clinicianID int) (int, error) {
		return counter.CountPatients(ctx, clinicianID)
	}
}

func removePatientFunc(remover patientRemover) RemovePatient {
	return func(ctx context.Context, p *Patient, clinicianID int) error {
		return remover.RemovePatient(ctx, p, clinicianID)
	}
}

func editPatientAddressFunc(editer patientAddressEditer) EditPatientAddress {
	return func(ctx context.Context, p *Patient, clinicianID int) error {
		return editer.EditPatientAddress(ctx, p, clinicianID)
	}
}
