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
	mockAccountGetter struct {
		account deiz.ClinicianAccount
		err     error
	}
	mockClinicianAccountCreater struct {
		err error
	}
	mockClinicianRegistrationVerifier struct {
		complete bool
		err      error
	}
	mockClinicianRegistrationCompleter struct {
		err error
	}
)

func (c *mockClinicianAccountCreater) CreateClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error {
	return c.err
}

func (c *mockAccountGetter) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return c.account, c.err
}
func (s *mockClinicianRegistrationVerifier) IsClinicianRegistrationComplete(ctx context.Context, email string) (bool, error) {
	return s.complete, s.err
}

func (s *mockClinicianRegistrationCompleter) CompleteClinicianRegistration(ctx context.Context, c *deiz.Clinician, password string, clinicianID int) error {
	return s.err
}

func TestEnsureClinicianRegistrationComplete(t *testing.T) {
	var tests = []struct {
		description string

		getter    *mockClinicianGetterByEmail
		verifier  *mockClinicianRegistrationVerifier
		completer *mockClinicianRegistrationCompleter

		inEmail    string
		inPassword string

		outError error
	}{
		{
			description: "should fail to validate email && password provided",
			outError:    deiz.ErrorStructValidation,
		},
		{
			description: "should not find user with given email",
			getter:      &mockClinicianGetterByEmail{err: errors.New("failed to get clinician by email")},
			inEmail:     "valid email",
			inPassword:  "valid password",
			outError:    errors.New("failed to get clinician by email"),
		},
		{
			description: "should fail to verify is registration is complete",
			getter:      &mockClinicianGetterByEmail{},
			verifier:    &mockClinicianRegistrationVerifier{err: errors.New("failed to access registration status")},
			inEmail:     "valid email",
			inPassword:  "valid password",
			outError:    errors.New("failed to access registration status"),
		},
		{
			description: "should confirm clinician is registered",
			getter:      &mockClinicianGetterByEmail{},
			verifier:    &mockClinicianRegistrationVerifier{complete: true},
			inEmail:     "valid email",
			inPassword:  "valid password",
		},
		{
			description: "should fail to complete registration",
			getter:      &mockClinicianGetterByEmail{},
			verifier:    &mockClinicianRegistrationVerifier{complete: false},
			completer:   &mockClinicianRegistrationCompleter{err: errors.New("failed to complete registration")},
			inEmail:     "valid email",
			inPassword:  "valid password",
			outError:    errors.New("failed to complete registration"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := account.Usecase{
				ClinicianGetterByEmail: test.getter,
				RegistrationVerifier:   test.verifier,
				RegistrationCompleter:  test.completer,
			}
			err := u.EnsureClinicianRegistrationComplete(context.Background(), test.inEmail, test.inPassword)
			assert.Equal(t, test.outError, err)
		})
	}
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
			r := account.Usecase{
				AccountGetter: test.accountGetter,
			}
			acc, err := r.GetClinicianAccount(context.Background(), test.inClinicianID)
			assert.Equal(t, test.outError, err)
			assert.Equal(t, test.outAccount, acc)
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
			r := account.Usecase{
				AccountCreater: test.creater,
			}
			err := r.AddClinicianAccount(context.Background(), test.inClinicianAccount)
			assert.Equal(t, test.outError, err)
		})
	}
}