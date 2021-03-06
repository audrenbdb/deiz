package deiz

import (
	"github.com/audrenbdb/deiz/valid"
	"strings"
)

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

func (p *Patient) FullName() string {
	return p.Surname + " " + p.Name
}

func (p *Patient) IsSet() bool {
	return p.ID != 0
}

func (p *Patient) IsNotSet() bool {
	return !p.IsSet()
}

func (p *Patient) IsEmailSet() bool {
	return p.Email != ""
}

func (p *Patient) IsValid() bool {
	return len(p.Name) >= 2 && len(p.Surname) >= 2 && valid.Phone(p.Phone) && (valid.Email(p.Email) || p.Email == "")
}

func (p *Patient) IsInvalid() bool {
	return !p.IsValid()
}

func (p *Patient) Sanitize() {
	p.Name = strings.TrimSpace(p.Name)
	p.Name = strings.ToUpper(p.Name)
	p.Surname = strings.TrimSpace(p.Surname)
	p.Surname = strings.Title(strings.ToLower(p.Surname))
	p.Email = strings.TrimSpace(p.Email)
	p.Email = strings.ToLower(p.Email)
	p.Phone = strings.TrimSpace(p.Phone)
	p.Phone = strings.Title(p.Phone)
}
