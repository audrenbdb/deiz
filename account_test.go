package deiz_test

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	mockAccountGetter struct {
		account deiz.ClinicianAccount
		err     error
	}
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
	mockClinicianAccountCreater struct {
		err error
	}
)

var validAddress = deiz.Address{
	ID:       1,
	Line:     "Test",
	PostCode: 10000,
	City:     "Test",
}

func (c *mockClinicianAccountCreater) CreateClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error {
	return c.err
}

func (c *mockAccountGetter) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return c.account, c.err
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

func TestGetClinicianAccount(t *testing.T) {
	var tests = []struct {
		description string

		accountGetter *mockAccountGetter

		inClinicianID int

		outAccount deiz.ClinicianAccount
		outError   error
	}{
		{
			description:   "should fail to get clinician account from the repo",
			accountGetter: &mockAccountGetter{err: errors.New("fail")},
			outError:      errors.New("fail"),
		},
		{
			description:   "should succeed to get clinician account",
			accountGetter: &mockAccountGetter{account: deiz.ClinicianAccount{Clinician: deiz.Clinician{ID: 1}}},
			outAccount:    deiz.ClinicianAccount{Clinician: deiz.Clinician{ID: 1}},
			outError:      nil,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := deiz.Repo{
				ClinicianAccount: deiz.ClinicianAccountRepo{
					AccountGetter: test.accountGetter,
				},
			}
			acc, err := r.GetClinicianAccount(context.Background(), test.inClinicianID)
			assert.Equal(t, test.outError, err)
			assert.Equal(t, test.outAccount, acc)
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
			r := deiz.Repo{
				ClinicianAccount: deiz.ClinicianAccountRepo{
					OfficeAddressCreater: test.adder,
				},
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
			r := deiz.Repo{
				ClinicianAccount: deiz.ClinicianAccountRepo{
					ClinicianAddressCreater: test.adder,
				},
			}
			err := r.AddClinicianAddress(context.Background(), test.inAddress, test.inClinicianID)
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
			r := deiz.Repo{
				ClinicianAccount: deiz.ClinicianAccountRepo{
					AddressUpdater:           test.addressUpdater,
					AddressOwnershipVerifier: test.verifier,
				},
			}
			err := r.UpdateClinicianAddress(context.Background(), test.inAddress, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}

func TestAddClinicianAccount(t *testing.T) {
	var tests = []struct {
		description string

		creater *mockClinicianAccountCreater

		inClinicianAccount *deiz.ClinicianAccount
		outError           error
	}{
		{
			description: "should fail to create clinician account",
			creater:     &mockClinicianAccountCreater{err: errors.New("failed to create")},
			outError:    errors.New("failed to create"),
		},
		{
			description: "should succeed",
			creater:     &mockClinicianAccountCreater{},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := deiz.Repo{
				ClinicianAccount: deiz.ClinicianAccountRepo{
					AccountCreater: test.creater,
				},
			}
			err := r.AddClinicianAccount(context.Background(), test.inClinicianAccount)
			assert.Equal(t, test.outError, err)
		})
	}
}
