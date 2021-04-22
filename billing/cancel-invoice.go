package billing

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func (c *CancelInvoiceUsecase) CancelInvoice(ctx context.Context, invoiceToCancel *deiz.BookingInvoice) error {
	invoiceToCancel.RemoveBooking()
	if err := setInvoiceIdentifier(ctx, invoiceToCancel, c.counter); err != nil {
		return err
	}
	return c.saver.SaveCorrectingBookingInvoice(ctx, invoiceToCancel)
}

type (
	correctiveInvoiceSaver interface {
		SaveCorrectingBookingInvoice(ctx context.Context, correctiveInvoice *deiz.BookingInvoice) error
	}
)

func NewCancelInvoiceUsecase(deps CancelInvoiceDeps) *CancelInvoiceUsecase {
	return &CancelInvoiceUsecase{
		counter: deps.Counter,
		saver:   deps.Saver,
	}
}

type CancelInvoiceDeps struct {
	Counter invoicesCounter
	Saver   correctiveInvoiceSaver
}

type CancelInvoiceUsecase struct {
	counter invoicesCounter
	saver   correctiveInvoiceSaver
}
