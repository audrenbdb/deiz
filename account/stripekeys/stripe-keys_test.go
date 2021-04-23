package stripekeys

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCrypter struct {
	bytes []byte
	err   error
}

type mockStripeKeysUpdater struct {
	err error
}

func (m *mockCrypter) StringToBytes(key string) ([]byte, error) {
	return m.bytes, m.err
}

func (m *mockStripeKeysUpdater) UpdateClinicianStripeKeys(ctx context.Context, pk string, sk []byte, clinicianID int) error {
	return m.err
}

func TestSetClinicianStripeKeys(t *testing.T) {
	validKey := "validpublickey"

	var tests = []struct {
		description string

		publicKey        string
		secretKeyInput   string
		clinicianIDInput int

		errorOutput error

		usecase Usecase
	}{
		{
			description: "should fail to validate key length",
			errorOutput: deiz.ErrorStructValidation,
		},
		{
			description: "should fail to encrypt key",

			publicKey:      validKey,
			secretKeyInput: validKey,
			errorOutput:    deiz.GenericError,

			usecase: Usecase{Crypter: &mockCrypter{err: deiz.GenericError}},
		},
		{
			description:    "should fail to update",
			publicKey:      validKey,
			secretKeyInput: validKey,
			errorOutput:    deiz.GenericError,

			usecase: Usecase{Crypter: &mockCrypter{}, StripeKeysUpdater: &mockStripeKeysUpdater{err: deiz.GenericError}},
		},
		{
			description:    "should succeed",
			publicKey:      validKey,
			secretKeyInput: validKey,
			usecase:        Usecase{Crypter: &mockCrypter{}, StripeKeysUpdater: &mockStripeKeysUpdater{}},
		},
	}

	for _, test := range tests {
		err := test.usecase.SetClinicianStripeKeys(context.Background(), test.publicKey, test.secretKeyInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
