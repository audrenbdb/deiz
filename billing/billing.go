package billing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

//invoice identifier format
const invoiceIDFormat = "DEIZ-%d-%08d"

//common interface for billing usecases
type (
	invoicesCounter interface {
		CountClinicianInvoices(ctx context.Context, clinicianID int) (int, error)
	}
	pdfInvoiceCreater interface {
		CreateBookingInvoicePDF(i *deiz.BookingInvoice) (*bytes.Buffer, error)
	}
	periodInvoicesGetter interface {
		GetPeriodBookingInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error)
	}
)

type mailInvoiceDeps struct {
	mailer     invoiceMailer
	pdfCreater pdfInvoiceCreater
	invoice    *deiz.BookingInvoice
	recipient  string
}

func mailInvoice(s mailInvoiceDeps) error {
	pdf, err := s.pdfCreater.CreateBookingInvoicePDF(s.invoice)
	if err != nil {
		return err
	}
	return s.mailer.MailBookingInvoice(s.invoice, pdf, s.recipient)
}

func setInvoiceIdentifier(ctx context.Context, invoice *deiz.BookingInvoice, counter invoicesCounter) error {
	count, err := counter.CountClinicianInvoices(ctx, invoice.ClinicianID)
	if err != nil {
		return err
	}
	invoice.Identifier = fmt.Sprintf(invoiceIDFormat, invoice.ClinicianID, count+1)
	return nil
}
