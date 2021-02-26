package deiz

import (
	"time"
)

type BookingInvoice struct {
	ID              int           `json:"id"`
	Booking         Booking       `json:"booking"`
	CreatedAt       time.Time     `json:"createdAt"`
	Identifier      string        `json:"identifier"`
	Sender          []string      `json:"sender"`
	Recipient       []string      `json:"recipient"`
	CityAndDate     string        `json:"cityAndDate"`
	DeliveryDate    time.Time     `json:"deliveryDate"`
	DeliveryDateStr string        `json:"deliveryDateStr"`
	Label           string        `json:"label"`
	PriceBeforeTax  int64         `json:"priceBeforeTax"`
	PriceAfterTax   int64         `json:"priceAfterTax"`
	TaxFee          float32       `json:"taxFee"`
	Exemption       string        `json:"exemption"`
	PaymentMethod   PaymentMethod `json:"paymentMethod"`
	Canceled        bool          `json:"canceled"`
}

type PaymentMethod struct {
	ID   int    `json:"id" validator:"required"`
	Name string `json:"name" validator:"required"`
}

/*
func seePeriodBookingInvoicesSummaryPDFFunc(tz clinicianTimezoneGetter, pdf periodBookingInvoicesSummaryPDFGetter, invoicesGetter periodBookingInvoicesGetter) SeePeriodInvoicesSummaryPDF {
	return func(ctx context.Context, start time.Time, end time.Time, clinicianID int) (*bytes.Buffer, error) {
		loc, err := getClinicianTimezoneLoc(ctx, clinicianID, tz)
		if err != nil {
			return nil, err
		}
		invoices, err := invoicesGetter.GetPeriodBookingInvoices(ctx, start, end, clinicianID)
		if err != nil {
			return nil, err
		}
		var totalBeforeTax int64
		var totalAfterTax int64
		for _, i := range invoices {
			totalBeforeTax = totalBeforeTax + i.PriceBeforeTax
			totalAfterTax = totalAfterTax + i.PriceAfterTax
		}
		return pdf.GetPeriodBookingInvoicesSummaryPDF(ctx, invoices, start, end, totalBeforeTax, totalAfterTax, loc)
	}
}

/*
func getCancelBookingURL(deleteID string) string {
	cancelURL, _ := url.Parse("https://deiz.fr")
	cancelURL.Path += "api/public/appointments/delete"
	params := url.Values{}
	params.Add("id", deleteID)
	cancelURL.RawQuery = params.Encode()
	return cancelURL.String()
}
*/
