package deiz

import "context"

type Address struct {
	ID       int    `json:"id"`
	Line     string `json:"line"`
	PostCode int    `json:"postCode"`
	City     string `json:"city"`
}

func (a *Address) isValid() bool {
	if len(a.Line) < 2 {
		return false
	}
	if a.PostCode < 10000 {
		return false
	}
	if len(a.City) < 2 {
		return false
	}
	return true
}

type (
	AddressUpdater interface {
		UpdateAddress(ctx context.Context, address *Address) error
	}
)
