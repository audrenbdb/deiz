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

func (f *ContactForm) Valid() bool {
	return f.ClinicianID != 0 && len(f.Name) >= 2 && len(f.Message) >= 2 && len(f.Email) >= 6
}

func (f *ContactForm) Invalid() bool {
	return !f.Valid()
}

func (f *GetInTouchForm) Valid() bool {
	return len(f.Email) >= 6 && len(f.Phone) >= 10 && len(f.Name) >= 2 && len(f.City) >= 2 && len(f.Job) >= 2
}

func (f *GetInTouchForm) Invalid() bool {
	return !f.Valid()
}
