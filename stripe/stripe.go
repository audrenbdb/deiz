package stripe

import (
	"context"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

const (
	successPaymentURL = "https://deiz.fr/payment/success/"
	failPaymentURL    = "https://deiz.fr/payment/cancel/"
)

type service struct{}

//CreateSession creates stripe session token that can be used by a client to open a stripe payment checkout form
func (s *service) CreateSession(ctx context.Context, amount int64, sk string) (string, error) {
	sc := &client.API{}
	sc.Init(sk, nil)
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("payment"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Entretien"),
					},
					UnitAmount: &amount,
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successPaymentURL),
		CancelURL:  stripe.String(failPaymentURL),
		Locale:     stripe.String("auto"),
	}
	session, err := sc.CheckoutSessions.New(params)
	return session.ID, err
}

func NewService() *service {
	return &service{}
}
