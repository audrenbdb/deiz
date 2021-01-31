package http

import (
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

func handleGetBookingsPendingPayment(getBookingsPendingPayments deiz.ListBookingsPendingPayment) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		b, err := getBookingsPendingPayments(ctx, clinicianID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}
