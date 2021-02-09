package account_test

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	mockClinicianOfficeAddressCreater struct {
		err error
	}
	mockClinicianAddressCreater struct {
		err error
	}
	mockAddressOwnershipVerifier struct {
		own bool
		err error
	}
	mockAddressUpdater struct {
		err error
	}
)

var validAddress = deiz.Address{
	ID:       1,
	Line:     "Test",
	PostCode: 10000,
	City:     "Test",
}

func (a *mockClinicianOfficeAddressCreater) CreateClinicianOfficeAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	return a.err
}

func (a *mockClinicianAddressCreater) CreateClinicianAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	return a.err
}

func (m *mockAddressOwnershipVerifier) IsAddressToClinician(ctx context.Context, a *deiz.Address, clinicianID int) (bool, error) {
	return m.own, m.err
}

func (c *mockAddressUpdater) UpdateAddress(ctx context.Context, address *deiz.Address) error {
	return c.err
}

func TestAddressIsValid(t *testing.T) {
	var tests = []struct {
		description string

		inAddress *deiz.Address

		isValid bool
	}{
		{
			description: "should fail with too low line",
			inAddress:   &deiz.Address{},
		},
		{
			description: "should fail with incorrect post code",
			inAddress:   &deiz.Address{Line: "test", PostCode: 0},
		},
		{
			description: "should succeed with correct address",
			inAddress:   &deiz.Address{Line: "test", PostCode: 10000, City: "TEST"},
			isValid:     true,
		},
		{
			description: "should return false with wrong city",
			inAddress:   &deiz.Address{Line: "test", PostCode: 10000},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			valid := account.IsAddressValid(test.inAddress)
			assert.Equal(t, test.isValid, valid)
		})
	}
}

func TestAddClinianOfficeAddress(t *testing.T) {
	var tests = []struct {
		description string

		adder *mockClinicianOfficeAddressCreater

		inAddress     *deiz.Address
		inClinicianID int

		outError error
	}{
		{
			description: "should fail to verify the address struct",
			inAddress:   &deiz.Address{PostCode: 0},
			outError:    deiz.ErrorStructValidation,
		},
		{
			description: "should fail to save office address in repo",
			adder:       &mockClinicianOfficeAddressCreater{err: errors.New("failed")},
			inAddress:   &deiz.Address{Line: "Test", PostCode: 10000, City: "Test"},
			outError:    errors.New("failed"),
		},
		{
			description: "should succeed in adding address",
			adder:       &mockClinicianOfficeAddressCreater{},
			inAddress:   &deiz.Address{Line: "Test", PostCode: 10000, City: "Test"},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := account.Usecase{
				OfficeAddressCreater: test.adder,
			}
			err := r.AddClinicianOfficeAddress(context.Background(), test.inAddress, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}
func TestAddClinianAddress(t *testing.T) {
	var tests = []struct {
		description string

		adder *mockClinicianAddressCreater

		inAddress     *deiz.Address
		inClinicianID int

		outError error
	}{
		{
			description: "should fail to verify the address struct",
			inAddress:   &deiz.Address{PostCode: 0},
			outError:    deiz.ErrorStructValidation,
		},
		{
			description: "should fail to save office address in repo",
			adder:       &mockClinicianAddressCreater{err: errors.New("failed")},
			inAddress:   &deiz.Address{Line: "Test", PostCode: 10000, City: "Test"},
			outError:    errors.New("failed"),
		},
		{
			description: "should succeed in adding address",
			adder:       &mockClinicianAddressCreater{},
			inAddress:   &deiz.Address{Line: "Test", PostCode: 10000, City: "Test"},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := account.Usecase{
				HomeAddressCreater: test.adder,
			}
			err := u.AddClinicianHomeAddress(context.Background(), test.inAddress, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}

func TestUpdateClinicianAddress(t *testing.T) {
	var tests = []struct {
		description string

		verifier       *mockAddressOwnershipVerifier
		addressUpdater *mockAddressUpdater

		inAddress     *deiz.Address
		inClinicianID int

		outError error
	}{
		{
			description: "should fail to validate the struct",
			inAddress:   &deiz.Address{},
			outError:    deiz.ErrorStructValidation,
		},
		{
			description: "should fail while trying to verify ownership",
			verifier:    &mockAddressOwnershipVerifier{err: errors.New("fail to verify")},
			inAddress:   &validAddress,
			outError:    errors.New("fail to verify"),
		},
		{
			description: "should fail to authorize ownership",
			verifier:    &mockAddressOwnershipVerifier{own: false},
			inAddress:   &validAddress,
			outError:    deiz.ErrorUnauthorized,
		},
		{
			description:    "should fail to update the address",
			verifier:       &mockAddressOwnershipVerifier{own: true},
			addressUpdater: &mockAddressUpdater{err: errors.New("fail to update")},
			inAddress:      &validAddress,
			outError:       errors.New("fail to update"),
		},
		{
			description:    "should succeed",
			verifier:       &mockAddressOwnershipVerifier{own: true},
			addressUpdater: &mockAddressUpdater{},
			inAddress:      &validAddress,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := account.Usecase{
				AddressUpdater:           test.addressUpdater,
				AddressOwnerShipVerifier: test.verifier,
			}
			err := u.UpdateClinicianAddress(context.Background(), test.inAddress, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}
