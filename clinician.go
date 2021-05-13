package deiz

type Clinician struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Profession string `json:"profession"`
	Adeli      Adeli  `json:"adeli"`
}

type Adeli struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
}

func (c *Clinician) FullName() string {
	return c.Surname + " " + c.Name
}
