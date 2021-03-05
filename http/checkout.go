package http

import (
	"context"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type (
	StripePaymentSessionCreater interface {
		CreateStripePaymentSession(ctx context.Context, amount int64, clinicianID int) (string, error)
	}
)

func handleGetSessionCheckout(creater StripePaymentSessionCreater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID, err := strconv.Atoi(c.QueryParam("clinicianId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		amount, err := strconv.Atoi(c.QueryParam("amount"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		session, err := creater.CreateStripePaymentSession(ctx, int64(amount), clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, session)
	}
}
