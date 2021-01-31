package deiz

import "context"

type (
	clinicianStripeSecretKeyGetter interface {
		GetClinicianStripeSecretKey(ctx context.Context, clinicianID int) ([]byte, error)
	}
	bytesDecrypter interface {
		DecryptBytes(ctx context.Context, b []byte) (string, error)
	}
	stripePaymentSessionCreater interface {
		CreateStripePaymentSession(ctx context.Context, amount int64, sk string) (string, error)
	}
)

type (
	CreateStripePaymentSession func(ctx context.Context, amount int64, clinicianID int) (string, error)
)

func creatStripePaymentSessionFunc(getter clinicianStripeSecretKeyGetter,
	decrypt bytesDecrypter,
	payment stripePaymentSessionCreater) CreateStripePaymentSession {
	return func(ctx context.Context, amount int64, clinicianID int) (string, error) {
		k, err := getter.GetClinicianStripeSecretKey(ctx, clinicianID)
		if err != nil {
			return "", nil
		}
		key, err := decrypt.DecryptBytes(ctx, k)
		if err != nil {
			return "", nil
		}
		return payment.CreateStripePaymentSession(ctx, amount, key)
	}
}
