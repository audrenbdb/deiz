package deiz

import "fmt"

type Address struct {
	ID       int    `json:"id"`
	Line     string `json:"line"`
	PostCode int    `json:"postCode"`
	City     string `json:"city"`
}

func (a *Address) IsSet() bool {
	return a.ID != 0
}

func (a *Address) IsNotSet() bool {
	return !a.IsSet()
}

func (a *Address) IsValid() bool {
	return len(a.Line) >= 2 && a.PostCode >= 10000 && len(a.City) >= 2
}

func (a *Address) IsInvalid() bool {
	return !a.IsValid()
}

func (a *Address) ToString() string {
	if a.ID == 0 {
		return ""
	}
	return fmt.Sprintf("%s, %d %s", a.Line, a.PostCode, a.City)
}
