package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
)

//CancelInvoice creates a corrective invoice to erase a previous one with wrong data
//French invoice system does not allow deletion of invoice, only correction.
func (c *CancelInvoiceUsecase) CancelInvoice(ctx context.Context, invoiceToCancel *deiz.BookingInvoice) error {
	invoiceToCancel.RemoveBooking()
	if err := setInvoiceIdentifier(ctx, invoiceToCancel, c.Counter); err != nil {
		return err
	}
	return c.Saver.SaveCorrectingBookingInvoice(ctx, invoiceToCancel)
}

type (
	correctiveInvoiceSaver interface {
		SaveCorrectingBookingInvoice(ctx context.Context, correctiveInvoice *deiz.BookingInvoice) error
	}
)

type CancelInvoiceUsecase struct {
	Counter invoicesCounter
	Saver   correctiveInvoiceSaver
}
