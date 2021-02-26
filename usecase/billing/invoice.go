package billing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	BookingInvoiceCreater interface {
		CreateBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, clinicianID int) error
	}
	BookingInvoicePDFCreater interface {
		CreateBookingInvoicePDF(i *deiz.BookingInvoice) (*bytes.Buffer, error)
	}
	InvoicesSummaryPDFCreater interface {
		CreateInvoicesSummaryPDF(i []deiz.BookingInvoice, start, end time.Time, loc *time.Location) (*bytes.Buffer, error)
	}
	ClinicianBoundChecker interface {
		IsPatientTiedToClinician(ctx context.Context, p *deiz.Patient, clinicianID int) (bool, error)
	}
	BookingInvoiceMailer interface {
		MailBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, invoicePDF *bytes.Buffer, sendTo string) error
	}
	InvoicesSummaryMailer interface {
		MailInvoicesSummary(ctx context.Context, summaryPDF *bytes.Buffer, sendTo string) error
	}
	InvoicesCounter interface {
		CountClinicianInvoices(ctx context.Context, clinicianID int) (int, error)
	}
	PeriodInvoicesGetter interface {
		GetPeriodBookingInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error)
	}
)

func (u *Usecase) GetPeriodInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error) {
	return u.PeriodInvoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
}

func (u *Usecase) GenerateBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, clinicianID int, sendToPatient bool) error {
	const invoiceIDFormat = "DEIZ-%d-%08d"
	bound, err := u.ClinicianBoundChecker.IsPatientTiedToClinician(ctx, &i.Booking.Patient, clinicianID)
	if err != nil {
		return err
	}
	if !bound {
		return deiz.ErrorUnauthorized
	}
	if i.TaxFee < 0 || i.PaymentMethod.ID <= 0 || i.PriceAfterTax < i.PriceBeforeTax || i.CityAndDate == "" || i.PriceAfterTax < 0 {
		return deiz.ErrorStructValidation
	}
	count, err := u.InvoicesCounter.CountClinicianInvoices(ctx, clinicianID)
	if err != nil {
		return err
	}
	i.Identifier = fmt.Sprintf(invoiceIDFormat, clinicianID, count+1)
	err = u.BookingInvoiceCreater.CreateBookingInvoice(ctx, i, clinicianID)
	if err != nil {
		return err
	}
	if sendToPatient {
		return u.MailBookingInvoice(ctx, i, i.Booking.Patient.Email)
	}
	return nil
}

func (u *Usecase) MailBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, sendTo string) error {
	invoicePDF, err := u.BookingInvoicePDFCreater.CreateBookingInvoicePDF(i)
	if err != nil {
		return err
	}
	return u.BookingInvoiceMailer.MailBookingInvoice(ctx, i, invoicePDF, sendTo)
}

func (u *Usecase) MailPeriodInvoicesSummary(ctx context.Context, start, end time.Time, tzName string, sendTo string, clinicianID int) error {
	invoices, err := u.PeriodInvoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
	if err != nil {
		return err
	}
	tz, err := time.LoadLocation(tzName)
	if err != nil {
		return err
	}
	pdf, err := u.InvoicesSummaryPDFCreater.CreateInvoicesSummaryPDF(invoices, start, end, tz)
	if err != nil {
		return err
	}
	return u.MailPeriodInvoicesSummary(ctx, pdf)
	return nil
}
