package deiz

//Patient uses the application to book clinician appointment
type Patient struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Phone   string  `json:"phone"`
	Email   string  `json:"email"`
	Note    string  `json:"note"`
	Address Address `json:"address"`
}
