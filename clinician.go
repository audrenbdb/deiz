package deiz

type Clinician struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	Address    Address `json:"address"`
	Profession string  `json:"profession"`
	Adeli      Adeli   `json:"adeli"`
}

type Adeli struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
}
