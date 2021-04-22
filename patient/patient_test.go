package patient_test

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/patient"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	mockPatientSearcher struct {
		patients []deiz.Patient
		err      error
	}
	mockPatientCreater struct {
		err error
	}
)

func (m *mockPatientCreater) CreatePatient(ctx context.Context, p *deiz.Patient, clinicianID int) error {
	return m.err
}

func (m *mockPatientSearcher) SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error) {
	return m.patients, m.err
}

func TestAddClinician(t *testing.T) {
	var tests = []struct {
		description string

		creater mockPatientCreater

		inPatient     deiz.Patient
		inClinicianID int
		outError      error
	}{
		{
			description: "should fail to validate",
			outError:    deiz.ErrorStructValidation,
		},
		{
			description: "should fail to create user",
			creater:     mockPatientCreater{err: errors.New("failed to create")},
			inPatient:   deiz.Patient{Name: "toto", Phone: "010101100101", Surname: "Hey", Email: "legit@legit.com"},
			outError:    errors.New("failed to create"),
		},
		{
			description: "should succeed",
			creater:     mockPatientCreater{},
			inPatient:   deiz.Patient{Name: "toto", Phone: "010101100101", Surname: "Hey", Email: "legit@legit.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := patient.Usecase{
				Creater: &test.creater,
			}
			err := u.AddPatient(context.Background(), &test.inPatient, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}

func TestIsValid(t *testing.T) {
	var tests = []struct {
		description string

		inPatient deiz.Patient
		outValid  bool
	}{
		{
			description: "should pass",
			inPatient:   deiz.Patient{Name: "toto", Phone: "010101100101", Surname: "Hey", Email: "legit@legit.com"},
			outValid:    true,
		},
		{
			description: "should fail invalid name",
			inPatient:   deiz.Patient{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			valid := patient.IsPatientValid(&test.inPatient)
			assert.Equal(t, test.outValid, valid)
		})
	}
}

func TestSearchPatient(t *testing.T) {
	var tests = []struct {
		description string

		searcher *mockPatientSearcher

		inSearch      string
		inClinicianID int

		outPatients []deiz.Patient
		outError    error
	}{
		{
			description: "should fail to search",
			searcher:    &mockPatientSearcher{err: errors.New("failed to search")},
			outError:    errors.New("failed to search"),
		},
		{
			description: "should pass",
			searcher:    &mockPatientSearcher{patients: []deiz.Patient{deiz.Patient{ID: 1}}},
			outPatients: []deiz.Patient{deiz.Patient{ID: 1}},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := patient.Usecase{
				Searcher: test.searcher,
			}
			p, err := u.SearchPatient(context.Background(), test.inSearch, test.inClinicianID)
			assert.Equal(t, test.outError, err)
			assert.Equal(t, test.outPatients, p)
		})
	}
}
