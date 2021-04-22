package billing

import (
	"bytes"
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockPeriodInvoicesGetter struct {
	invoices []deiz.BookingInvoice
	err      error
}

func (m *mockPeriodInvoicesGetter) GetPeriodBookingInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error) {
	return m.invoices, m.err
}

type mockInvoicesSummaryPDFCreater struct {
	err error
}

func (m *mockInvoicesSummaryPDFCreater) CreateInvoicesSummaryPDF(i []deiz.BookingInvoice, start, end time.Time) (*bytes.Buffer, error) {
	return nil, m.err
}

type mockInvoicesSummaryMailer struct {
	err error
}

func (m *mockInvoicesSummaryMailer) MailInvoicesSummary(summaryPDF *bytes.Buffer, start, end time.Time, sendTo string) error {
	return m.err
}

func TestMailInvoice(t *testing.T) {
	var tests = []struct {
		description string

		invoiceInput   *deiz.BookingInvoice
		recipientInput string
		errorOutput    error

		usecase MailInvoiceUsecase
	}{
		{
			description: "should fail to send email",

			errorOutput: deiz.GenericError,

			usecase: MailInvoiceUsecase{
				pdfInvoiceCreater: &mockPDFCreater{},
				invoiceMailer:     &mockInvoiceSender{err: deiz.GenericError},
			},
		},
		{
			description: "should fail to create pdf",

			errorOutput: deiz.GenericError,

			usecase: MailInvoiceUsecase{
				pdfInvoiceCreater: &mockPDFCreater{err: deiz.GenericError},
				invoiceMailer:     &mockInvoiceSender{},
			},
		},
		{
			description: "should succeed",

			usecase: MailInvoiceUsecase{
				pdfInvoiceCreater: &mockPDFCreater{},
				invoiceMailer:     &mockInvoiceSender{},
			},
		},
	}
	for _, test := range tests {
		err := test.usecase.MailInvoice(&validInvoice, test.recipientInput)
		assert.Equal(t, test.errorOutput, err)
	}
}

func TestMailInvoicesSummary(t *testing.T) {
	var tests = []struct {
		description string

		startInput       time.Time
		endInput         time.Time
		recipientInput   string
		clinicianIDInput int
		errorOutput      error

		usecase MailInvoiceUsecase
	}{
		{
			description: "should fail to get invoices summary",
			errorOutput: deiz.GenericError,

			usecase: MailInvoiceUsecase{
				invoicesGetter: &mockPeriodInvoicesGetter{err: deiz.GenericError},
			},
		},
		{
			description: "should fail to create summary pdf",
			errorOutput: deiz.GenericError,

			usecase: MailInvoiceUsecase{
				invoicesGetter:            &mockPeriodInvoicesGetter{},
				pdfInvoicesSummaryCreater: &mockInvoicesSummaryPDFCreater{err: deiz.GenericError},
			},
		},
		{
			description: "should fail to send summary pdf through email",
			errorOutput: deiz.GenericError,

			usecase: MailInvoiceUsecase{
				invoicesGetter:            &mockPeriodInvoicesGetter{},
				pdfInvoicesSummaryCreater: &mockInvoicesSummaryPDFCreater{},
				invoicesSummaryMailer:     &mockInvoicesSummaryMailer{err: deiz.GenericError},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.MailInvoicesSummary(context.Background(), test.startInput, test.endInput, test.recipientInput, test.clinicianIDInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
