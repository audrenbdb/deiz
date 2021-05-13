package account

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockGetClinicianAccount struct {
	account deiz.ClinicianAccount
	err     error
}

func (m *mockGetClinicianAccount) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return m.account, m.err
}

func TestGetClinicianAccountData(t *testing.T) {
	var tests = []struct {
		description string

		credInput        deiz.Credentials
		clinicianIDInput int
		accountOuput     deiz.ClinicianAccount
		errorOutput      error

		usecase GetDataUsecase
	}{
		{
			description: "Should fail to retrieve clinician account",

			credInput:   deiz.Credentials{Role: deiz.ClinicianRole},
			errorOutput: deiz.GenericError,
			usecase: GetDataUsecase{
				&mockGetClinicianAccount{err: deiz.GenericError},
			},
		},
		{
			description: "Should filter out private booking motive if public",

			credInput: deiz.Credentials{Role: deiz.PatientRole},
			accountOuput: deiz.ClinicianAccount{BookingMotives: []deiz.BookingMotive{
				{Public: true, ID: 1},
			}},
			usecase: GetDataUsecase{
				&mockGetClinicianAccount{account: deiz.ClinicianAccount{
					Business: deiz.Business{ID: 1},
					BookingMotives: []deiz.BookingMotive{
						{Public: false, ID: 1}, {Public: true, ID: 1},
					}}},
			},
		},
	}

	for _, test := range tests {
		acc, err := test.usecase.GetClinicianAccountData(context.Background(), test.clinicianIDInput, test.credInput)

		assert.Equal(t, test.errorOutput, err)
		assert.Equal(t, test.accountOuput, acc)
	}
}

/*
func TestGetClinicianAccountPublicData(t *testing.T) {
	var tests = []struct {
		description string

		clinicianIDInput int
		accountOuput     deiz.ClinicianAccountPublicData
		errorOutput      error

		usecase GetDataUsecase
	}{
		{
			description: "should fail to get account data",

			errorOutput: deiz.GenericError,

			usecase: GetDataUsecase{
				&mockGetClinicianAccount{err: deiz.GenericError},
			},
		},
		{
			description: "should split data into public data only",
			accountOuput: deiz.ClinicianAccountPublicData{
				PublicMotives: []deiz.BookingMotive{},
				RemoteAllowed: true,
			},
			usecase: GetDataUsecase{
				&mockGetClinicianAccount{account: deiz.ClinicianAccount{
					BookingMotives:   []deiz.BookingMotive{},
					CalendarSettings: deiz.CalendarSettings{RemoteAllowed: true},
				}},
			},
		},
	}
	for _, test := range tests {
		acc, err := test.usecase.GetClinicianAccountPublicData(context.Background(), test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
		assert.Equal(t, test.accountOuput, acc)
	}
}
*/
