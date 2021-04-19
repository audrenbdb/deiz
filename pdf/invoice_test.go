package pdf

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGenerateBookingInvoicePDF(t *testing.T) {
	pdf := NewService(
		"oxygen",
		"oxygen.ttf",
		filepath.Join("../", "assets", "fonts"))

	invoice := &deiz.BookingInvoice{
		Identifier:     "test",
		Sender:         []string{"test sender", "test"},
		Recipient:      []string{"test recipient", "recipient"},
		CityAndDate:    "TestCity and date",
		DeliveryDate:   "01/01/0001",
		Label:          "test",
		PriceBeforeTax: 5000,
		PriceAfterTax:  6000,
		TaxFee:         20.0,
		Exemption:      "290",
		PaymentMethod:  deiz.PaymentMethod{Name: "Carte bancaire"},
	}

	buffer, err := pdf.GenerateBookingInvoicePDF(context.Background(), invoice)
	assert.NoError(t, err)

	invoiceBytes := buffer.Bytes()
	err = ioutil.WriteFile("test_single_invoice.Pdf", invoiceBytes, 0644)
	assert.NoError(t, err)
}

func TestGetPeriodInvoicesSummary(t *testing.T) {
	pdf := NewService(
		"oxygen",
		"oxygen.ttf",
		filepath.Join("../", "assets", "fonts"))
	invoice := deiz.BookingInvoice{
		Identifier:     "test",
		Sender:         []string{"test sender", "test"},
		Recipient:      []string{"test recipient", "recipient"},
		CityAndDate:    "TestCity and date",
		DeliveryDate:   "01/01/0001",
		Label:          "test",
		PriceBeforeTax: 5000,
		PriceAfterTax:  6000,
		TaxFee:         20.0,
		Exemption:      "290",
		PaymentMethod:  deiz.PaymentMethod{Name: "Carte bancaire"},
	}
	invoices := []deiz.BookingInvoice{invoice, invoice}
	buffer, err := pdf.GetPeriodInvoicesSummary(invoices, "01/01/0001", "01/02/0001", 5000, 5000)
	assert.NoError(t, err)

	invoiceBytes := buffer.Bytes()
	err = ioutil.WriteFile("test_summary.Pdf", invoiceBytes, 0644)
	assert.NoError(t, err)
}
