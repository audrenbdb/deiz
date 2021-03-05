package deiz

type ClinicianAccount struct {
	Clinician        Clinician        `json:"clinician"`
	Business         Business         `json:"business"`
	OfficeAddresses  []Address        `json:"officeAddresses"`
	StripePublicKey  string           `json:"stripePublicKey"`
	OfficeHours      []OfficeHours    `json:"officeHours"`
	BookingMotives   []BookingMotive  `json:"bookingMotives"`
	CalendarSettings CalendarSettings `json:"calendarSettings"`
}

type ClinicianAccountPublicData struct {
	Clinician       Clinician       `json:"clinician"`
	StripePublicKey string          `json:"stripePublicKey"`
	PublicMotives   []BookingMotive `json:"bookingMotives"`
	ClinicianTz     string          `json:"clinicianTz"`
	RemoteAllowed   bool            `json:"remoteAllowed"`
}
