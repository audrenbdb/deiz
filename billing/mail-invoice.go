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
		mailer:     m.invoiceMailer,
		pdfCreater: m.pdfInvoiceCreater,
		recipient:  recipient,
	})
}

func (m *MailInvoiceUsecase) MailInvoicesSummary(ctx context.Context, start, end time.Time, recipient string, clinicianID int) error {
	invoices, err := m.invoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
	if err != nil {
		return err
	}
	invoicesPDF, err := m.pdfInvoicesSummaryCreater.CreateInvoicesSummaryPDF(invoices, start, end)
	if err != nil {
		return err
	}
	return m.invoicesSummaryMailer.MailInvoicesSummary(invoicesPDF, start, end, recipient)
}

type (
	invoicesSummaryPDFCreater interface {
		CreateInvoicesSummaryPDF(i []deiz.BookingInvoice, start, end time.Time) (*bytes.Buffer, error)
	}
	invoicesSummaryMailer interface {
		MailInvoicesSummary(summaryPDF *bytes.Buffer, start, end time.Time, sendTo string) error
	}
)

func NewMailInvoiceUsecase(deps MailInvoiceDeps) *MailInvoiceUsecase {
	return &MailInvoiceUsecase{
		invoiceMailer:             deps.InvoiceMailer,
		pdfInvoiceCreater:         deps.PdfInvoiceCreater,
		pdfInvoicesSummaryCreater: deps.PdfInvoicesSummaryCreater,
		invoicesGetter:            deps.InvoicesGetter,
		invoicesSummaryMailer:     deps.InvoicesSummaryMailer,
	}
}

type MailInvoiceDeps struct {
	InvoiceMailer             invoiceMailer
	InvoicesSummaryMailer     invoicesSummaryMailer
	PdfInvoicesSummaryCreater invoicesSummaryPDFCreater
	PdfInvoiceCreater         pdfInvoiceCreater
	InvoicesGetter            periodInvoicesGetter
}

type MailInvoiceUsecase struct {
	invoiceMailer             invoiceMailer
	pdfInvoiceCreater         pdfInvoiceCreater
	pdfInvoicesSummaryCreater invoicesSummaryPDFCreater
	invoicesGetter            periodInvoicesGetter
	invoicesSummaryMailer     invoicesSummaryMailer
}
