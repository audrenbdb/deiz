package deiz

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"
)

const invoiceIDFormat = "DEIZ-%d-%08d"

type BookingInvoice struct {
	ID              int           `json:"id" validator:"required"`
	Booking         Booking       `json:"booking" validator:"required"`
	CreatedAt       time.Time     `json:"createdAt"`
	Identifier      string        `json:"identifier" validator:"required"`
	Sender          []string      `json:"sender" validator:"required"`
	Recipient       []string      `json:"recipient" validator:"required"`
	CityAndDate     string        `json:"cityAndDate" validator:"required"`
	DeliveryDate    time.Time     `json:"deliveryDate" validator:"required"`
	DeliveryDateStr string        `json:"bookingDateStr" validator:"required"`
	Label           string        `json:"label" validator:"required"`
	PriceBeforeTax  int64         `json:"amount" validator:"required"`
	PriceAfterTax   int64         `json:"amount" validator:"required"`
	TaxFee          float32       `json:"taxFee" validator:"min=0"`
	Exemption       string        `json:"exemption"`
	PaymentMethod   PaymentMethod `json:"paymentMethod" validator:"required"`
	Canceled        bool          `json:"canceled"`
}

type PaymentMethod struct {
	ID   int    `json:"id" validator:"required"`
	Name string `json:"name" validator:"required"`
}

//External drivers to call
type (
	bookingInvoiceCreater interface {
		CreateBookingInvoice(ctx context.Context, invoice *BookingInvoice, clinicianID int) error
	}
	clinicianInvoicesCounter interface {
		CountClinicianInvoices(ctx context.Context, clinicianID int) (int, error)
	}
	bookingsPendingPaymentGetter interface {
		GetBookingsPendingPayment(ctx context.Context, clinicianID int) ([]Booking, error)
	}
	bookingInvoiceMailer interface {
		MailBookingInvoice(ctx context.Context, invoice *BookingInvoice, invoicePDF *bytes.Buffer) error
	}
	bookingInvoicePDFGenerater interface {
		GenerateBookingInvoicePDF(ctx context.Context, invoice *BookingInvoice) (*bytes.Buffer, error)
	}
	periodBookingInvoicesSummaryPDFGetter interface {
		GetPeriodBookingInvoicesSummaryPDF(ctx context.Context, invoices []BookingInvoice, start, end time.Time, totalBeforeTax, totalAfterTax int64, clinicianTz *time.Location) (*bytes.Buffer, error)
	}
	periodBookingInvoicesGetter interface {
		GetPeriodBookingInvoices(ctx context.Context, start time.Time, end time.Time, clinicianID int) ([]BookingInvoice, error)
	}
)

//core functions
type (
	//CreateBookingInvoice creates a new invoice for a given appointment
	CreateBookingInvoice func(ctx context.Context, i *BookingInvoice, clinicianID int) error
	//ListBookingsPendingPayment lists all prior to today appointments that have not been paid yet
	ListBookingsPendingPayment func(ctx context.Context, clinicianID int) ([]Booking, error)
	//MailBookingInvoice sends an email with an invoice attached as pdf
	MailBookingInvoice func(ctx context.Context, i *BookingInvoice) error
	//SeeInvoicePDF returns a bytes buffer to be decoded as PDF by the client
	SeeInvoicePDF func(ctx context.Context, invoice *BookingInvoice) (*bytes.Buffer, error)
	//SeePeriodInvoicesSummary returns a pdf containing all the bill over a period
	//With a total earnings over that given period
	SeePeriodInvoicesSummaryPDF func(ctx context.Context, start time.Time, end time.Time, clinicianID int) (*bytes.Buffer, error)
)

func createBookingInvoiceFunc(creater bookingInvoiceCreater,
	counter clinicianInvoicesCounter) CreateBookingInvoice {
	return func(ctx context.Context, i *BookingInvoice, clinicianID int) error {
		count, err := counter.CountClinicianInvoices(ctx, clinicianID)
		if err != nil {
			return err
		}
		i.Identifier = fmt.Sprintf(invoiceIDFormat, clinicianID, count+1)
		return creater.CreateBookingInvoice(ctx, i, clinicianID)
	}
}

/*
func seePeriodBookingInvoicesSummaryPDFFunc(tz clinicianTimezoneGetter, pdf periodBookingInvoicesSummaryPDFGetter, invoicesGetter periodBookingInvoicesGetter) SeePeriodInvoicesSummaryPDF {
	return func(ctx context.Context, start time.Time, end time.Time, clinicianID int) (*bytes.Buffer, error) {
		loc, err := getClinicianTimezoneLoc(ctx, clinicianID, tz)
		if err != nil {
			return nil, err
		}
		invoices, err := invoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
		if err != nil {
			return nil, err
		}
		var totalBeforeTax int64
		var totalAfterTax int64
		for _, i := range invoices {
			totalBeforeTax = totalBeforeTax + i.PriceBeforeTax
			totalAfterTax = totalAfterTax + i.PriceAfterTax
		}
		return pdf.GetPeriodBookingInvoicesSummaryPDF(ctx, invoices, start, end, totalBeforeTax, totalAfterTax, loc)
	}
}

*/

func listBookingsPendingPaymentFunc(lister bookingsPendingPaymentGetter) ListBookingsPendingPayment {
	return func(ctx context.Context, clinicianID int) ([]Booking, error) {
		return lister.GetBookingsPendingPayment(ctx, clinicianID)
	}
}

func mailBookingInvoiceFunc(generator bookingInvoicePDFGenerater,
	mailer bookingInvoiceMailer) MailBookingInvoice {
	return func(ctx context.Context, i *BookingInvoice) error {
		invoicePDF, err := generator.GenerateBookingInvoicePDF(ctx, i)
		if err != nil {
			return err
		}
		return mailer.MailBookingInvoice(ctx, i, invoicePDF)
	}
}

func seeInvoicePDFFunc(generater bookingInvoicePDFGenerater) SeeInvoicePDF {
	return func(ctx context.Context, invoice *BookingInvoice) (*bytes.Buffer, error) {
		return generater.GenerateBookingInvoicePDF(ctx, invoice)
	}
}

func getCancelBookingURL(deleteID string) string {
	cancelURL, _ := url.Parse("https://deiz.fr")
	cancelURL.Path += "api/public/appointments/delete"
	params := url.Values{}
	params.Add("id", deleteID)
	cancelURL.RawQuery = params.Encode()
	return cancelURL.String()
}
