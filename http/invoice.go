package http

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type (
	BookingInvoiceGenerater interface {
		GenerateBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, clinicianID int, sendToPatient bool) error
	}
	BookingInvoiceMailer interface {
		MailBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, sendTo string) error
	}
	PeriodInvoicesSummaryMailer interface {
		MailPeriodInvoicesSummary(ctx context.Context, start, end time.Time, tzName string, sendTo string, clinicianID int) error
	}
	PeriodInvoicesGetter interface {
		GetPeriodInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error)
	}
)

func handlePostPDFBookingInvoicesPeriodSummary(mailer PeriodInvoicesSummaryMailer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		type post struct {
			SendTo   string `json:"sendTo"`
			Timezone string `json:"timezone"`
		}
		var p post
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		startParam, err := strconv.ParseInt(c.QueryParam("start"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		endParam, err := strconv.ParseInt(c.QueryParam("end"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return mailer.MailPeriodInvoicesSummary(ctx, time.Unix(startParam, 0), time.Unix(endParam, 0), p.Timezone, p.SendTo, clinicianID)
	}
}

func handlePostPDFBookingInvoice(mailer BookingInvoiceMailer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		type post struct {
			SendTo  string              `json:"sendTo"`
			Invoice deiz.BookingInvoice `json:"invoice"`
		}
		var p post
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		err := mailer.MailBookingInvoice(ctx, &p.Invoice, p.SendTo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostBookingInvoice(generater BookingInvoiceGenerater) echo.HandlerFunc {
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
		err = generater.GenerateBookingInvoice(ctx, &i, clinicianID, sendToPatient)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, i)
	}
}

func handleGetPeriodInvoices(getter PeriodInvoicesGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		startParam, err := strconv.ParseInt(c.QueryParam("start"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		endParam, err := strconv.ParseInt(c.QueryParam("end"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		invoices, err := getter.GetPeriodInvoices(ctx, time.Unix(startParam, 0), time.Unix(endParam, 0), clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, invoices)
	}
}

/*
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
*/
