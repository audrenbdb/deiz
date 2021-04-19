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
		CreateInvoicesSummaryPDF(i []deiz.BookingInvoice, start, end time.Time) (*bytes.Buffer, error)
	}
	ClinicianBoundChecker interface {
		IsPatientTiedToClinician(ctx context.Context, p *deiz.Patient, clinicianID int) (bool, error)
	}
	BookingInvoiceMailer interface {
		MailBookingInvoice(i *deiz.BookingInvoice, invoicePDF *bytes.Buffer, sendTo string) error
	}
	InvoicesSummaryMailer interface {
		MailInvoicesSummary(summaryPDF *bytes.Buffer, start, end time.Time, sendTo string) error
	}
	InvoicesCounter interface {
		CountClinicianInvoices(ctx context.Context, clinicianID int) (int, error)
	}
	PeriodInvoicesGetter interface {
		GetPeriodBookingInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error)
	}
	BookingInvoiceCanceler interface {
		CancelBookingInvoice(ctx context.Context, originalInvoiceID int, i *deiz.BookingInvoice, clinicianID int) error
	}
)

func (u *Usecase) GetPeriodInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error) {
	return u.PeriodInvoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
}

func (u *Usecase) GenerateBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, clinicianID int, sendToPatient bool) error {
	bound, err := u.ClinicianBoundChecker.IsPatientTiedToClinician(ctx, &i.Booking.Patient, clinicianID)
	if err != nil {
		return err
	}
	if !bound {
		return deiz.ErrorUnauthorized
	}
	if i.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	err = u.generateInvoiceIdentifier(ctx, i, clinicianID)
	if err != nil {
		return err
	}
	err = u.BookingInvoiceCreater.CreateBookingInvoice(ctx, i, clinicianID)
	if err != nil {
		return err
	}
	if sendToPatient {
		return u.MailBookingInvoice(ctx, i, i.Booking.Patient.Email)
	}
	return nil
}

func (u *Usecase) generateInvoiceIdentifier(ctx context.Context, i *deiz.BookingInvoice, clinicianID int) error {
	count, err := u.InvoicesCounter.CountClinicianInvoices(ctx, clinicianID)
	if err != nil {
		return fmt.Errorf("unable to generate invoice identifier: %s", err)
	}
	i.SetIdentifier(clinicianID, count)
	return nil
}

func (u *Usecase) CancelBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, clinicianID int) error {
	i.RemoveBooking()
	if i.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	err := u.generateInvoiceIdentifier(ctx, i, clinicianID)
	if err != nil {
		return err
	}
	return u.BookingInvoiceCanceler.CancelBookingInvoice(ctx, i.ID, i, clinicianID)
}

func (u *Usecase) MailBookingInvoice(ctx context.Context, i *deiz.BookingInvoice, sendTo string) error {
	invoicePDF, err := u.BookingInvoicePDFCreater.CreateBookingInvoicePDF(i)
	if err != nil {
		return err
	}
	return u.BookingInvoiceMailer.MailBookingInvoice(i, invoicePDF, sendTo)
}

func (u *Usecase) MailPeriodInvoicesSummary(ctx context.Context, start, end time.Time, sendTo string, clinicianID int) error {
	invoices, err := u.PeriodInvoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
	if err != nil {
		return err
	}
	pdf, err := u.InvoicesSummaryPDFCreater.CreateInvoicesSummaryPDF(invoices, start, end)
	if err != nil {
		return err
	}
	return u.InvoicesSummaryMailer.MailInvoicesSummary(pdf, start, end, sendTo)
}
