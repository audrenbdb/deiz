package address

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockAddressDeleter struct {
	err error
}

func (m *mockAddressDeleter) DeleteAddress(ctx context.Context, addressID int) error {
	return m.err
}

func TestDeleteAddress(t *testing.T) {
	var tests = []struct {
		description string

		addressIDInput   int
		clinicianIDInput int
		errorOutput      error

		usecase DeleteAddressUsecase
	}{
		{
			description:    "Should fail to fail to verify address ownership",
			addressIDInput: 1,
			errorOutput:    deiz.ErrorUnauthorized,

			usecase: DeleteAddressUsecase{
				AccountGetter: &mockAccountGetter{account: deiz.ClinicianAccount{}},
			},
		},
		{
			description:    "should fail to delete address",
			addressIDInput: 1,
			errorOutput:    deiz.GenericError,

			usecase: DeleteAddressUsecase{
				AccountGetter: &mockAccountGetter{account: deiz.ClinicianAccount{
					OfficeAddresses: []deiz.Address{deiz.Address{ID: 1}},
				}},
				AddressDeleter: &mockAddressDeleter{err: deiz.GenericError},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.DeleteAddress(context.Background(), test.addressIDInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
