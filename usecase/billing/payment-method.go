package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	PaymentMethodsGetter interface {
		GetPaymentMethods(ctx context.Context) ([]deiz.PaymentMethod, error)
	}
)

func (u *Usecase) GetAvailablePaymentMethods(ctx context.Context) ([]deiz.PaymentMethod, error) {
	return u.PaymentMethodsGetter.GetPaymentMethods(ctx)
}
