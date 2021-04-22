package billing

import (
	"bytes"
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockInvoiceSaver struct {
	err error
}

func (m *mockInvoiceSaver) SaveBookingInvoice(ctx context.Context, i *deiz.BookingInvoice) error {
	return m.err
}

type mockPDFCreater struct {
	err error
}

func (m *mockPDFCreater) CreateBookingInvoicePDF(i *deiz.BookingInvoice) (*bytes.Buffer, error) {
	return nil, m.err
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		usecase            CreateInvoiceUsecase
		description        string
		invoiceInput       *deiz.BookingInvoice
		sendToPatientInput bool

		errorOutput error
	}{
		{
			description: "should fail to validate empty struct",

			invoiceInput: &invalidInvoice,
			errorOutput:  deiz.ErrorStructValidation,

			usecase: CreateInvoiceUsecase{},
		},
		{
			description: "should fail to set identifier because count fail",

			invoiceInput: &validInvoice,
			errorOutput:  deiz.GenericError,

			usecase: CreateInvoiceUsecase{
				Counter: &mockInvoicesCounter{err: deiz.GenericError},
			},
		},
		{
			description: "should fail to save invoice",

			invoiceInput: &validInvoice,
			errorOutput:  deiz.GenericError,

			usecase: CreateInvoiceUsecase{
				Counter: &mockInvoicesCounter{},
				Saver:   &mockInvoiceSaver{err: deiz.GenericError},
			},
		},
		{
			description: "should fail to create invoice pdf",

			sendToPatientInput: true,
			invoiceInput:       &validInvoice,

			errorOutput: deiz.GenericError,

			usecase: CreateInvoiceUsecase{
				Counter:    &mockInvoicesCounter{},
				Saver:      &mockInvoiceSaver{},
				PdfCreater: &mockPDFCreater{err: deiz.GenericError},
			},
		},
		{
			description: "should fail to send invoice pdf",

			sendToPatientInput: true,
			invoiceInput:       &validInvoice,

			errorOutput: deiz.GenericError,

			usecase: CreateInvoiceUsecase{
				Counter:    &mockInvoicesCounter{},
				Saver:      &mockInvoiceSaver{},
				PdfCreater: &mockPDFCreater{},
				Mailer:     &mockInvoiceSender{err: deiz.GenericError},
			},
		},
	}

	for _, test := range tests {
		ctx := context.Background()
		err := test.usecase.CreateInvoice(ctx, test.invoiceInput, test.sendToPatientInput)
		assert.Equal(t, test.errorOutput, err)
	}
}
