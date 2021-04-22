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
	if err := setInvoiceIdentifier(ctx, invoice, i.counter); err != nil {
		return err
	}
	if err := i.saver.SaveBookingInvoice(ctx, invoice); err != nil {
		return err
	}
	if sendToPatient {
		return i.send(invoice)
	}
	return nil
}

func (i *CreateInvoiceUsecase) send(invoice *deiz.BookingInvoice) error {
	return mailInvoice(mailInvoiceDeps{
		pdfCreater: i.pdfCreater,
		mailer:     i.mailer,
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
	counter    invoicesCounter
	saver      invoiceSaver
	mailer     invoiceMailer
	pdfCreater pdfInvoiceCreater
}

type CreateInvoiceDeps struct {
	Counter    invoicesCounter
	Saver      invoiceSaver
	Mailer     invoiceMailer
	PdfCreater pdfInvoiceCreater
}

func NewCreateInvoiceUsecase(deps CreateInvoiceDeps) *CreateInvoiceUsecase {
	return &CreateInvoiceUsecase{
		counter:    deps.Counter,
		saver:      deps.Saver,
		mailer:     deps.Mailer,
		pdfCreater: deps.PdfCreater,
	}
}
