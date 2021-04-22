package stripekeys

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type Usecase struct {
	Crypter           crypter
	StripeKeysUpdater updater
}

type (
	updater interface {
		UpdateClinicianStripeKeys(ctx context.Context, pk string, sk []byte, clinicianID int) error
	}
	crypter interface {
		StringToBytes(key string) ([]byte, error)
	}
)

func (u *Usecase) SetClinicianStripeKeys(ctx context.Context, pk, sk string, clinicianID int) error {
	if len(pk) < 7 || len(sk) < 7 {
		return deiz.ErrorStructValidation
	}
	encryptedSecretKey, err := u.Crypter.StringToBytes(sk)
	if err != nil {
		return err
	}
	return u.StripeKeysUpdater.UpdateClinicianStripeKeys(ctx, pk, encryptedSecretKey, clinicianID)
}
