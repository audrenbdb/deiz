package billing

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
)

type (
	BytesDecrypter interface {
		DecryptBytes(strBytes []byte) (string, error)
	}
	ClinicianStripeSecretKeyGetter interface {
		GetClinicianStripeSecretKey(ctx context.Context, clinicianID int) ([]byte, error)
	}
	StripePaymentSessionCreater interface {
		CreateStripePaymentSession(ctx context.Context, amount int64, sk string) (string, error)
	}
)

func (u *Usecase) CreateStripePaymentSession(ctx context.Context, amount int64, clinicianID int) (string, error) {
	if amount <= 0 {
		return "", deiz.ErrorUnauthorized
	}
	k, err := u.ClinicianStripeSecretKeyGetter.GetClinicianStripeSecretKey(ctx, clinicianID)
	if err != nil {
		return "", err
	}
	key, err := u.BytesDecrypter.DecryptBytes(k)
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", errors.New("no key provided")
	}
	return u.StripePaymentSessionCreater.CreateStripePaymentSession(ctx, amount, key)
}
