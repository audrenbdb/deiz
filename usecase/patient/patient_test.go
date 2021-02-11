package patient_test

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase/patient"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	mockPatientSearcher struct {
		patients []deiz.Patient
		err      error
	}
)

func (m *mockPatientSearcher) SearchPatient(ctx context.Context, search string, clinicianID int) ([]deiz.Patient, error) {
	return m.patients, m.err
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
