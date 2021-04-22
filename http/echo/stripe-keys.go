package echo

import (
	"context"
	"github.com/labstack/echo"
	"net/http"
)

type (
	stripeKeysSetter interface {
		SetClinicianStripeKeys(ctx context.Context, pk, sk string, clinicianID int) error
	}
)

func handlePatchStripeKeys(setter stripeKeysSetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		type keys struct {
			PublicKey string `json:"publicKey"`
			SecretKey string `json:"secretKey"`
		}
		var k keys
		if err := c.Bind(&k); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err := setter.SetClinicianStripeKeys(ctx, k.PublicKey, k.SecretKey, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
