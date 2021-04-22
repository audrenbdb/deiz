package clinician

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockUpdater struct {
	err error
}

func (m *mockUpdater) UpdateClinicianPhone(ctx context.Context, phone string, clinicianID int) error {
	return m.err
}

func (m *mockUpdater) UpdateClinicianEmail(ctx context.Context, email string, clinicianID int) error {
	return m.err
}

func (m *mockUpdater) UpdateClinicianAddress(ctx context.Context, address deiz.Address, clinicianID int) error {
	return m.err
}

func TestEditClinicianPhone(t *testing.T) {
	validPhone := "05505050505"
	var tests = []struct {
		description string

		phoneInput       string
		clinicianIDInput int
		errorOutput      error

		usecase EditUsecase
	}{
		{
			description: "should fail to validate phone",

			errorOutput: deiz.ErrorStructValidation,

			usecase: EditUsecase{},
		},
		{
			description: "should fail to update phone",

			phoneInput:  validPhone,
			errorOutput: deiz.GenericError,

			usecase: EditUsecase{
				PhoneUpdater: &mockUpdater{err: deiz.GenericError},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.EditlinicianPhone(context.Background(), test.phoneInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
	}
}

func TestEditClinicianEmail(t *testing.T) {
	validEmail := "test@test.com"
	var tests = []struct {
		description string

		emailInput       string
		clinicianIDInput int
		errorOutput      error

		usecase EditUsecase
	}{
		{
			description: "should fail to validate email",

			errorOutput: deiz.ErrorStructValidation,

			usecase: EditUsecase{},
		},
		{
			description: "should fail to update email",

			emailInput:  validEmail,
			errorOutput: deiz.GenericError,

			usecase: EditUsecase{
				EmailUpdater: &mockUpdater{err: deiz.GenericError},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.EditClinicianEmail(context.Background(), test.emailInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
