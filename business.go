package deiz

//Business relates to professional business informations
type Business struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Identifier   string       `json:"identifier"`
	TaxExemption TaxExemption `json:"taxExemption"`
}

type TaxExemption struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}
