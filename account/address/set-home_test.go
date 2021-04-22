package address

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockHomeAddressManager struct {
	err error
}

func (m *mockHomeAddressManager) CreateClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	return m.err
}

func (m *mockHomeAddressManager) SetClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error {
	return m.err
}

func TestSetHomeAddress(t *testing.T) {
	validNewAddress := deiz.Address{
		ID:       0,
		PostCode: 29000,
		City:     "TEST",
		Line:     "TEST",
	}

	validExistingAddress := deiz.Address{
		ID:       1,
		PostCode: 29000,
		City:     "TEST",
		Line:     "TEST",
	}

	var tests = []struct {
		description string

		addressInput     *deiz.Address
		clinicianIDInput int
		errorOutput      error

		usecase SetHomeUsecase
	}{
		{
			description: "should fail to validate the address",

			addressInput: &deiz.Address{},
			errorOutput:  deiz.ErrorStructValidation,

			usecase: SetHomeUsecase{},
		},
		{
			description:  "should attempt to CREATE an address and fail",
			addressInput: &validNewAddress,
			errorOutput:  deiz.GenericError,
			usecase: SetHomeUsecase{
				manager: &mockHomeAddressManager{
					err: deiz.GenericError,
				},
			},
		},
		{
			description:  "should attempt to UPDATE an address and fail",
			addressInput: &validExistingAddress,
			errorOutput:  deiz.GenericError,
			usecase: SetHomeUsecase{
				manager: &mockHomeAddressManager{
					err: deiz.GenericError,
				},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.SetHomeAddress(context.Background(), test.addressInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
