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

//SetClinicianStripeKeys saves encrypted stripe keys of a given clinician
func (u *Usecase) SetClinicianStripeKeys(ctx context.Context, publicKey, secretKey string, clinicianID int) error {
	if keysInvalid(publicKey, secretKey) {
		return deiz.ErrorStructValidation
	}
	encryptedSecretKey, err := u.Crypter.StringToBytes(secretKey)
	if err != nil {
		return err
	}
	return u.StripeKeysUpdater.UpdateClinicianStripeKeys(ctx, publicKey, encryptedSecretKey, clinicianID)
}

func keysInvalid(publicKey, secretKey string) bool {
	return len(publicKey) < 7 || len(secretKey) < 7
}
