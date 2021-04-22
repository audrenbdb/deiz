package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCorrectingInvoiceSaver struct {
	err error
}

func (m *mockCorrectingInvoiceSaver) SaveCorrectingBookingInvoice(ctx context.Context, correctiveInvoice *deiz.BookingInvoice) error {
	return m.err
}
func TestCancelInvoice(t *testing.T) {
	var tests = []struct {
		description string

		invoiceInput *deiz.BookingInvoice
		errorOutput  error

		usecase CancelInvoiceUsecase
	}{
		{
			description: "Should remove the booking when canceling an invoice",

			invoiceInput: &deiz.BookingInvoice{Booking: deiz.Booking{ID: 1}},

			usecase: CancelInvoiceUsecase{
				Counter: &mockInvoicesCounter{},
				Saver:   &mockCorrectingInvoiceSaver{},
			},
		},
		{
			description: "Should fail to generate an invoice identifier",

			invoiceInput: &validInvoice,
			errorOutput:  deiz.GenericError,

			usecase: CancelInvoiceUsecase{
				Counter: &mockInvoicesCounter{err: deiz.GenericError},
			},
		},
		{
			description: "should fail to save correcting invoice",

			invoiceInput: &validInvoice,
			errorOutput:  deiz.GenericError,

			usecase: CancelInvoiceUsecase{
				Counter: &mockInvoicesCounter{},
				Saver:   &mockCorrectingInvoiceSaver{err: deiz.GenericError},
			},
		},
		{
			description: "should pass",

			invoiceInput: &validInvoice,
			usecase: CancelInvoiceUsecase{
				Counter: &mockInvoicesCounter{},
				Saver:   &mockCorrectingInvoiceSaver{},
			},
		},
	}

	for _, test := range tests {
		err := test.usecase.CancelInvoice(context.Background(), test.invoiceInput)
		assert.Equal(t, test.errorOutput, err)
		assert.True(t, test.invoiceInput.Booking.ID == 0)
	}
}
