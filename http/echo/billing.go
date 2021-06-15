package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func handleGetUnpaidBookings(getter usecase.UnpaidBookingsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		bookings, err := getter.GetUnpaidBookings(ctx, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handleGetSessionCheckout(creater usecase.StripeSessionCreater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID, err := getURLIntegerQueryParam(c, "clinicianId")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		amount, err := getURLIntegerQueryParam(c, "amount")
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

func handlePostPDFBookingInvoicesPeriodSummary(mailer usecase.InvoiceMailer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handlePostPDFBookingInvoice(mailer usecase.InvoiceMailer) echo.HandlerFunc {
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

func handlePostBookingInvoice(creater usecase.InvoiceCreater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		sendToPatient, err := strconv.ParseBool(c.QueryParam("sendToPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		invoice, err := getInvoiceFromRequest(c)
		invoice.ClinicianID = getCredFromEchoCtx(c).UserID
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

func handlePostCancelInvoice(canceler usecase.InvoiceCanceler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		invoice, err := getInvoiceFromRequest(c)
		invoice.ClinicianID = getCredFromEchoCtx(c).UserID
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := canceler.CancelInvoice(ctx, invoice); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())

		}
		return c.JSON(http.StatusOK, invoice)
	}
}

func handleGetPeriodInvoices(getter usecase.InvoicesGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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
	i.Booking.Clinician.ID = getCredFromEchoCtx(c).UserID
	return &i, nil
}
