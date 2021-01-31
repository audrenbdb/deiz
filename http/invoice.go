package http

import (
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func handlePostBookingInvoice(createInvoice deiz.CreateBookingInvoice, mailInvoice deiz.MailBookingInvoice, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		sendToPatient, err := strconv.ParseBool(c.QueryParam("sendToPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var i deiz.BookingInvoice
		if err := c.Bind(&i); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if err := validate.StructExceptCtx(ctx, i, "ID", "Booking"); err != nil {
			return c.JSON(http.StatusBadRequest, errValidating)
		}
		err = createInvoice(ctx, &i, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if sendToPatient {
			err = mailInvoice(ctx, &i)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		return c.JSON(http.StatusOK, i)
	}
}

func handleGetBookingInvoicePDF(getInvoicePDF deiz.SeeInvoicePDF, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var i deiz.BookingInvoice
		if err := c.Bind(&i); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if err := validate.StructCtx(ctx, i); err != nil {
			return c.JSON(http.StatusBadRequest, errValidating)
		}
		pdfBytes, err := getInvoicePDF(ctx, &i)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.Blob(http.StatusOK, "application/pdf", pdfBytes.Bytes())
	}
}
