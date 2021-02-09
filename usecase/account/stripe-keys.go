package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianStripeKeysUpdater interface {
		UpdateClinicianStripeKeys(ctx context.Context, pk string, sk []byte, clinicianID int) error
	}
)

func (u *Usecase) SetClinicianStripeKeys(ctx context.Context, pk, sk string, clinicianID int) error {
	if len(pk) < 7 || len(sk) < 7 {
		return deiz.ErrorStructValidation
	}
	skCrypted, err := u.StringToBytesCrypter.CryptStringToBytes(sk)
	if err != nil {
		return err
	}
	return u.StripeKeysUpdater.UpdateClinicianStripeKeys(ctx, pk, skCrypted, clinicianID)
}
