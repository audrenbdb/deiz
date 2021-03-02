package deiz

type Address struct {
	ID       int    `json:"id"`
	Line     string `json:"line"`
	PostCode int    `json:"postCode"`
	City     string `json:"city"`
}
