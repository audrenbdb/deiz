package deiz

type ContactForm struct {
	ClinicianID int    `json:"clinicianId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Message     string `json:"message"`
}

type GetInTouchForm struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Job   string `json:"job"`
	City  string `json:"city"`
}
