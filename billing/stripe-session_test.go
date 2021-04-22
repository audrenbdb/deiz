package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockStripeKeyGetter struct {
	err error
}

func (m *mockStripeKeyGetter) GetClinicianStripeSecretKey(ctx context.Context, clinicianID int) ([]byte, error) {
	return nil, m.err
}

func TestCreateStripePaymentSession(t *testing.T) {
	var tests = []struct {
		description string

		amountInput      int64
		clinicianIDInput int

		sessionOutput string
		errorOutput   error

		usecase CreateStripeSessionUsecase
	}{
		{
			description: "should fail because wrong amount input",

			amountInput: -5,
			errorOutput: deiz.ErrorStructValidation,

			usecase: CreateStripeSessionUsecase{},
		},
		{
			description: "should fail to get stripe key of the clinician",

			amountInput: 1,
			errorOutput: deiz.GenericError,

			usecase: CreateStripeSessionUsecase{
				secretKeyGetter: &mockStripeKeyGetter{err: deiz.GenericError},
			},
		},
	}

	for _, test := range tests {
		session, err := test.usecase.CreateStripePaymentSession(
			context.Background(), test.amountInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
		assert.Equal(t, test.sessionOutput, session)
	}
}
