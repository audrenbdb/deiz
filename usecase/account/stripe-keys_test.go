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
	mockStripeKeysUpdater struct {
		err error
	}
)

func (s *mockStripeKeysUpdater) UpdateClinicianStripeKeys(ctx context.Context, pk string, sk []byte, clinicianID int) error {
	return s.err
}

func TestSetClinicianStripeKeys(t *testing.T) {
	var tests = []struct {
		description string

		crypter *mockCrypter
		updater *mockStripeKeysUpdater

		inPublicKey   string
		inSecretKey   string
		inClinicianID int

		outError error
	}{
		{
			description: "should fail to validate keys length",
			outError:    deiz.ErrorStructValidation,
		},
		{
			description: "should fail to encrypt secret key",
			crypter:     &mockCrypter{err: errors.New("fail to encode key")},
			inPublicKey: "valid public key",
			inSecretKey: "valid secret key",
			outError:    errors.New("fail to encode key"),
		},
		{
			description: "should fail to update keys in db",
			crypter:     &mockCrypter{},
			updater:     &mockStripeKeysUpdater{err: errors.New("fail to update keys")},
			inPublicKey: "valid public key",
			inSecretKey: "valid secret key",
			outError:    errors.New("fail to update keys"),
		},
		{
			description: "should pass and update with success",
			crypter:     &mockCrypter{},
			updater:     &mockStripeKeysUpdater{},
			inPublicKey: "valid public key",
			inSecretKey: "valid secret key",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := account.Usecase{
				StringToBytesCrypter: test.crypter,
				StripeKeysUpdater:    test.updater,
			}
			err := u.SetClinicianStripeKeys(context.Background(), test.inPublicKey, test.inSecretKey, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}
