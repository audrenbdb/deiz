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
	ID   int    `json:"id"`
	Name string `json:"name"`
}
