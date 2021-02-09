package deiz

//Business relates to professional business informations
type Business struct {
	ID           int          `json:"id" validate:"required"`
	Name         string       `json:"name" validate:"required"`
	Identifier   string       `json:"identifier" validate:"required"`
	TaxExemption TaxExemption `json:"taxExemption" validate:"required"`
}

type TaxExemption struct {
	ID   int    `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
}
