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
	return u.stripeSessionCreater.CreateSession(ctx, amount, key)
}

func (u *CreateStripeSessionUsecase) getDecryptedStripeKey(ctx context.Context, clinicianID int) (string, error) {
	k, err := u.secretKeyGetter.GetClinicianStripeSecretKey(ctx, clinicianID)
	if err != nil {
		return "", err
	}
	return decryptKey(u.crypter, k)
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
	crypter              crypter
	secretKeyGetter      stripeSecretKeyGetter
	stripeSessionCreater stripeSessionCreater
}

type CreateStripeSessionDeps struct {
	Crypter         crypter
	SecretKeyGetter stripeSecretKeyGetter
	SessionCreater  stripeSessionCreater
}

func NewCreateStripeSessionUsecase(deps CreateStripeSessionDeps) *CreateStripeSessionUsecase {
	return &CreateStripeSessionUsecase{
		crypter:              deps.Crypter,
		secretKeyGetter:      deps.SecretKeyGetter,
		stripeSessionCreater: deps.SessionCreater,
	}
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
