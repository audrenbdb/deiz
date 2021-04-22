package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type (
	invoiceCreater interface {
		CreateInvoice(ctx context.Context, invoice *deiz.BookingInvoice, sendToPatient bool) error
	}
	invoiceCanceler interface {
		CancelInvoice(ctx context.Context, invoice *deiz.BookingInvoice) error
	}
	invoiceMailer interface {
		MailInvoice(invoice *deiz.BookingInvoice, recipient string) error
		MailInvoicesSummary(ctx context.Context, start, end time.Time, recipient string, clinicianID int) error
	}
	invoicesGetter interface {
		GetPeriodInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error)
	}
	stripeSessionCreater interface {
		CreateStripePaymentSession(ctx context.Context, amount int64, clinicianID int) (string, error)
	}
	unpaidBookingsGetter interface {
		GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error)
	}
)

func handleGetUnpaidBookings(getter unpaidBookingsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		bookings, err := getter.GetUnpaidBookings(ctx, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handleGetSessionCheckout(creater stripeSessionCreater) echo.HandlerFunc {
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

func handlePostPDFBookingInvoicesPeriodSummary(mailer invoiceMailer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		type post struct {
			SendTo string `json:"sendTo"`
		}
		var p post
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		start, err := getTimeFromParam(c, "start")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		end, err := getTimeFromParam(c, "end")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return mailer.MailInvoicesSummary(ctx, start, end, p.SendTo, clinicianID)
	}
}

func handlePostPDFBookingInvoice(mailer invoiceMailer) echo.HandlerFunc {
	return func(c echo.Context) error {
		type post struct {
			SendTo  string              `json:"sendTo"`
			Invoice deiz.BookingInvoice `json:"invoice"`
		}
		var p post
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		err := mailer.MailInvoice(&p.Invoice, p.SendTo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostBookingInvoice(creater invoiceCreater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		sendToPatient, err := strconv.ParseBool(c.QueryParam("sendToPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		invoice, err := getInvoiceFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = creater.CreateInvoice(ctx, invoice, sendToPatient)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, invoice)
	}
}

func handlePostCancelInvoice(canceler invoiceCanceler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		invoice, err := getInvoiceFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := canceler.CancelInvoice(ctx, invoice); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())

		}
		return c.JSON(http.StatusOK, invoice)
	}
}

func handleGetPeriodInvoices(getter invoicesGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		start, err := getTimeFromParam(c, "start")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		end, err := getTimeFromParam(c, "end")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		invoices, err := getter.GetPeriodInvoices(ctx, start, end, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, invoices)
	}
}

func getInvoiceFromRequest(c echo.Context) (*deiz.BookingInvoice, error) {
	var i deiz.BookingInvoice
	if err := c.Bind(&i); err != nil {
		return nil, err
	}
	i.Booking.Clinician.ID = getCredFromEchoCtx(c).userID
	return &i, nil
}
