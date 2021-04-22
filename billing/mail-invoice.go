package billing

import (
	"bytes"
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

func (m *MailInvoiceUsecase) MailInvoice(invoice *deiz.BookingInvoice, recipient string) error {
	return mailInvoice(mailInvoiceDeps{
		invoice:    invoice,
		mailer:     m.InvoiceMailer,
		pdfCreater: m.PdfInvoiceCreater,
		recipient:  recipient,
	})
}

func (m *MailInvoiceUsecase) MailInvoicesSummary(ctx context.Context, start, end time.Time, recipient string, clinicianID int) error {
	invoices, err := m.InvoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
	if err != nil {
		return err
	}
	invoicesPDF, err := m.PdfInvoicesSummaryCreater.CreateInvoicesSummaryPDF(invoices, start, end)
	if err != nil {
		return err
	}
	return m.InvoicesSummaryMailer.MailInvoicesSummary(invoicesPDF, start, end, recipient)
}

type (
	invoicesSummaryPDFCreater interface {
		CreateInvoicesSummaryPDF(i []deiz.BookingInvoice, start, end time.Time) (*bytes.Buffer, error)
	}
	invoicesSummaryMailer interface {
		MailInvoicesSummary(summaryPDF *bytes.Buffer, start, end time.Time, sendTo string) error
	}
)

type MailInvoiceUsecase struct {
	InvoiceMailer             invoiceMailer
	PdfInvoiceCreater         pdfInvoiceCreater
	PdfInvoicesSummaryCreater invoicesSummaryPDFCreater
	InvoicesGetter            periodInvoicesGetter
	InvoicesSummaryMailer     invoicesSummaryMailer
}
