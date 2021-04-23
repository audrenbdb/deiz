package deiz

type ClinicianAccount struct {
	Clinician        Clinician        `json:"clinician"`
	Business         Business         `json:"business"`
	OfficeAddresses  []Address        `json:"officeAddresses"`
	StripePublicKey  string           `json:"stripePublicKey"`
	OfficeHours      []OfficeHours    `json:"officeHours"`
	BookingMotives   []BookingMotive  `json:"bookingMotives"`
	CalendarSettings CalendarSettings `json:"calendarSettings"`
	PaymentMethods   []PaymentMethod  `json:"paymentMethods"`
	TaxExemptions    []TaxExemption   `json:"taxExemptions"`
}

type LoginData struct {
	Email    string
	Password string
}

func (c *LoginData) IsInvalid() bool {
	return len(c.Email) < 4 || len(c.Password) < 6
}
