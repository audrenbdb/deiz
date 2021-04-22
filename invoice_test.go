package deiz_test

import (
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

var validInvoice = deiz.BookingInvoice{
	TaxFee:        1,
	PaymentMethod: deiz.PaymentMethod{ID: 1, Name: "A"},
	CityAndDate:   "A",
}

var invalidInvoice = deiz.BookingInvoice{}

func TestInvoiceValid(t *testing.T) {

	var tests = []struct {
		invoice deiz.BookingInvoice

		valid bool
	}{
		{
			invoice: validInvoice,
			valid:   true,
		},
		{
			invoice: invalidInvoice,
			valid:   false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.valid, test.invoice.IsValid())
	}
}
