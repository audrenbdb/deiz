package address

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockAddressCreater struct {
	err error
}

func (m *mockAddressCreater) CreateClinicianOfficeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	return m.err
}

func TestAddClinicianOfficeAddress(t *testing.T) {

	validNewAddress := deiz.Address{
		PostCode: 29000,
		City:     "TEST",
		Line:     "TEST",
	}

	var tests = []struct {
		description string

		addressInput     *deiz.Address
		clinicianIDInput int
		errorOutput      error

		usecase AddAddressUsecase
	}{
		{
			description: "should fail to validate address",

			addressInput: &deiz.Address{},
			errorOutput:  deiz.ErrorStructValidation,
		},
		{
			description: "should fail to create the address",

			addressInput: &validNewAddress,
			errorOutput:  deiz.GenericError,

			usecase: AddAddressUsecase{
				AddressCreater: &mockAddressCreater{err: deiz.GenericError},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.AddClinicianOfficeAddress(context.Background(), test.addressInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
