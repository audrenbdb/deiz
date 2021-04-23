package account

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockClinicianGetterByEmail struct {
	clinician deiz.Clinician
	err       error
}

func (m *mockClinicianGetterByEmail) GetClinicianByEmail(ctx context.Context, email string) (deiz.Clinician, error) {
	return m.clinician, m.err
}

type mockAuthChecker struct {
	enabled bool
	err     error
}

func (m *mockAuthChecker) IsClinicianAuthenticationEnabled(ctx context.Context, email string) (bool, error) {
	return m.enabled, m.err
}

type mockAuthEnabler struct {
	err error
}

func (m *mockAuthEnabler) EnableClinicianAuthentication(ctx context.Context, clinician *deiz.Clinician, password string) error {
	return m.err
}

func TestAllowLogin(t *testing.T) {
	validCredentials := deiz.LoginData{
		Email:    "random legit email",
		Password: "a random password",
	}

	var tests = []struct {
		description string

		credentialsInput deiz.LoginData

		errorOuput error

		usecase AllowLoginUsecase
	}{
		{
			description: "should fail to validate credential length",

			errorOuput: deiz.ErrorStructValidation,

			usecase: AllowLoginUsecase{},
		},
		{
			description:      "should fail to get clinician with given email",
			credentialsInput: validCredentials,
			errorOuput:       deiz.GenericError,

			usecase: AllowLoginUsecase{
				ClinicianGetter: &mockClinicianGetterByEmail{err: deiz.GenericError},
			},
		},
		{
			description:      "should fail to check if authentication is enabled",
			credentialsInput: validCredentials,
			errorOuput:       deiz.GenericError,

			usecase: AllowLoginUsecase{
				ClinicianGetter: &mockClinicianGetterByEmail{},
				AuthChecker:     &mockAuthChecker{err: deiz.GenericError},
			},
		},
		{
			description:      "should fail to enable authentication",
			credentialsInput: validCredentials,
			errorOuput:       deiz.GenericError,

			usecase: AllowLoginUsecase{
				ClinicianGetter: &mockClinicianGetterByEmail{},
				AuthChecker:     &mockAuthChecker{},
				AuthEnabler:     &mockAuthEnabler{err: deiz.GenericError},
			},
		},
		{
			description:      "should succeed if authentication is already enabled",
			credentialsInput: validCredentials,

			usecase: AllowLoginUsecase{
				ClinicianGetter: &mockClinicianGetterByEmail{},
				AuthChecker:     &mockAuthChecker{enabled: true},
			},
		},
		{
			description:      "should succeed in enabling non enabled to authenticate account",
			credentialsInput: validCredentials,

			usecase: AllowLoginUsecase{
				ClinicianGetter: &mockClinicianGetterByEmail{},
				AuthChecker:     &mockAuthChecker{enabled: false},
				AuthEnabler:     &mockAuthEnabler{},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.AllowLogin(context.Background(), test.credentialsInput)
		assert.Equal(t, test.errorOuput, err)
	}
}
