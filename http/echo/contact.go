package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	ContactFormToClinicianSender interface {
		SendContactFormToClinician(ctx context.Context, f deiz.ContactForm) error
	}
	GetInTouchSender interface {
		SendGetInTouchForm(ctx context.Context, f deiz.GetInTouchForm) error
	}
)

func handlePostContactFormToClinician(sender ContactFormToClinicianSender) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var f deiz.ContactForm
		if err := c.Bind(&f); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := sender.SendContactFormToClinician(ctx, f)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostGetInTouchForm(sender GetInTouchSender) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var f deiz.GetInTouchForm
		if err := c.Bind(&f); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := sender.SendGetInTouchForm(ctx, f)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
