package billing

import (
	"bytes"
	"context"
	"github.com/audrenbdb/deiz"
)

func (i *CreateInvoiceUsecase) CreateInvoice(ctx context.Context, invoice *deiz.BookingInvoice, sendToPatient bool) error {
	if invoice.IsInvalid() {
		return deiz.ErrorStructValidation
	}
	if err := setInvoiceIdentifier(ctx, invoice, i.Counter); err != nil {
		return err
	}
	if err := i.Saver.SaveBookingInvoice(ctx, invoice); err != nil {
		return err
	}
	if sendToPatient && invoice.Booking.Patient.IsEmailSet() {
		return i.send(invoice)
	}
	return nil
}

func (i *CreateInvoiceUsecase) send(invoice *deiz.BookingInvoice) error {
	return mailInvoice(mailInvoiceDeps{
		pdfCreater: i.PdfCreater,
		mailer:     i.Mailer,
		invoice:    invoice,
		recipient:  invoice.Booking.Patient.Email,
	})
}

type (
	invoiceSaver interface {
		SaveBookingInvoice(ctx context.Context, i *deiz.BookingInvoice) error
	}
	invoiceMailer interface {
		MailBookingInvoice(invoice *deiz.BookingInvoice, invoicePDF *bytes.Buffer, recipient string) error
	}
)

type CreateInvoiceUsecase struct {
	Counter    invoicesCounter
	Saver      invoiceSaver
	Mailer     invoiceMailer
	PdfCreater pdfInvoiceCreater
}
