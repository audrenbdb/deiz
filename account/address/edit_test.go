package address

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockAddressUpdater struct {
	err error
}

func (m *mockAddressUpdater) UpdateAddress(ctx context.Context, address *deiz.Address) error {
	return m.err
}

type mockAccountGetter struct {
	account deiz.ClinicianAccount
	err     error
}

func (m *mockAccountGetter) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return m.account, m.err
}

func TestEditAddress(t *testing.T) {
	validAddress := deiz.Address{
		ID:       1,
		PostCode: 29000,
		City:     "TEST",
		Line:     "TEST",
	}
	var tests = []struct {
		description string

		addressInput *deiz.Address
		credInput    deiz.Credentials
		errorOutput  error

		usecase EditAddressUsecase
	}{
		{
			description: "Shoud fail to validate address",

			addressInput: &deiz.Address{},
			errorOutput:  deiz.ErrorStructValidation,

			usecase: EditAddressUsecase{},
		},
		{
			description: "Should fail to get current account in order to get clinician addresses",

			addressInput: &validAddress,
			errorOutput:  deiz.GenericError,

			usecase: EditAddressUsecase{AccountGetter: &mockAccountGetter{err: deiz.GenericError}},
		},
		{
			description: "should fail to find address in clinician's address",

			addressInput: &validAddress,
			errorOutput:  deiz.ErrorUnauthorized,

			usecase: EditAddressUsecase{AccountGetter: &mockAccountGetter{account: deiz.ClinicianAccount{
				Clinician: deiz.Clinician{
					Address: deiz.Address{},
				},
				OfficeAddresses: nil,
			}}},
		},
		{
			description: "should fail to update",

			addressInput: &validAddress,
			errorOutput:  deiz.GenericError,

			usecase: EditAddressUsecase{AccountGetter: &mockAccountGetter{account: deiz.ClinicianAccount{
				Clinician: deiz.Clinician{
					Address: validAddress,
				},
				OfficeAddresses: nil,
			}},
				AddressUpdater: &mockAddressUpdater{err: deiz.GenericError}},
		},
	}

	for _, test := range tests {
		err := test.usecase.EditAddress(context.Background(), test.addressInput, test.credInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
