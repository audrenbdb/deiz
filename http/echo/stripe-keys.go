package echo

import (
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

func handlePatchStripeKeys(setter usecase.StripeKeysSetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

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
