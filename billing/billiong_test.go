package billing

import (
	"bytes"
	"context"
	"github.com/audrenbdb/deiz"
)

var validInvoice = deiz.BookingInvoice{
	TaxFee:        1,
	PaymentMethod: deiz.PaymentMethod{ID: 1, Name: "A"},
	CityAndDate:   "A",
}

var invalidInvoice = deiz.BookingInvoice{}

type mockInvoicesCounter struct {
	count int
	err   error
}

func (m *mockInvoicesCounter) CountClinicianInvoices(ctx context.Context, clinicianID int) (int, error) {
	return m.count, m.err
}

type mockInvoiceSender struct {
	err error
}

func (m *mockInvoiceSender) MailBookingInvoice(invoice *deiz.BookingInvoice, invoicePDF *bytes.Buffer, sendTo string) error {
	return m.err
}
