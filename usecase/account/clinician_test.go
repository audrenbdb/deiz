package account_test

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	mockClinicianGetterByEmail struct {
		clinician deiz.Clinician
		err       error
	}
)

func (s *mockClinicianGetterByEmail) GetClinicianByEmail(ctx context.Context, email string) (deiz.Clinician, error) {
	return s.clinician, s.err
}
