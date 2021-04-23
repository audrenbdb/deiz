package usecase

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	PatientUsecases struct {
		Searcher       PatientSearcher
		Adder          PatientAdder
		Editer         PatientEditer
		AddressAdder   PatientAddressAdder
		AddressEditer  PatientAddressEditer
		BookingsGetter PatientBookingsGetter
	}
)

type (
	PatientSearcher interface {
		SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error)
	}
	PatientAdder interface {
		AddPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	PatientEditer interface {
		EditPatient(ctx context.Context, p *deiz.Patient, clinicianID int) error
	}
	PatientAddressAdder interface {
		AddPatientAddress(ctx context.Context, a *deiz.Address, patientID int, clinicianID int) error
	}
	PatientAddressEditer interface {
		EditPatientAddress(ctx context.Context, a *deiz.Address, patientID int, clinicianID int) error
	}
	PatientBookingsGetter interface {
		GetPatientBookings(ctx context.Context, clinicianID, patientID int) ([]deiz.Booking, error)
	}
)
