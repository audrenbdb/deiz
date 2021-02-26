package http

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

type (
	PaymentMethodsGetter interface {
		GetAvailablePaymentMethods(ctx context.Context) ([]deiz.PaymentMethod, error)
	}
)

func handleGetPaymentMethods(getter PaymentMethodsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		methods, err := getter.GetAvailablePaymentMethods(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, methods)
	}
}
