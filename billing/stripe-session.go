package billing

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
)

func (u *CreateStripeSessionUsecase) CreateStripePaymentSession(ctx context.Context, amount int64, clinicianID int) (string, error) {
	if amount <= 0 {
		return "", deiz.ErrorStructValidation
	}
	key, err := u.getDecryptedStripeKey(ctx, clinicianID)
	if err != nil {
		return "", err
	}
	return u.StripeSessionCreater.CreateSession(ctx, amount, key)
}

func (u *CreateStripeSessionUsecase) getDecryptedStripeKey(ctx context.Context, clinicianID int) (string, error) {
	k, err := u.SecretKeyGetter.GetClinicianStripeSecretKey(ctx, clinicianID)
	if err != nil {
		return "", err
	}
	return decryptKey(u.Crypter, k)
}

func decryptKey(crypter crypter, k []byte) (string, error) {
	key, err := crypter.BytesToString(k)
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", errors.New("no key provided")
	}
	return key, nil
}

type CreateStripeSessionUsecase struct {
	Crypter              crypter
	SecretKeyGetter      stripeSecretKeyGetter
	StripeSessionCreater stripeSessionCreater
}

type (
	crypter interface {
		BytesToString(strBytes []byte) (string, error)
	}
	stripeSecretKeyGetter interface {
		GetClinicianStripeSecretKey(ctx context.Context, clinicianID int) ([]byte, error)
	}
	stripeSessionCreater interface {
		CreateSession(ctx context.Context, amount int64, sk string) (string, error)
	}
)
