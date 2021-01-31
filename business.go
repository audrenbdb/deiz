package deiz

import "context"

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

type (
	clinicianBusinessEditer interface {
		EditClinicianBusiness(ctx context.Context, business *Business, clinicianID int) error
	}
)

type (
	EditClinicianBusiness func(ctx context.Context, business *Business, clinicianID int) error
)

func editClinicianBusinessFunc(edit clinicianBusinessEditer) EditClinicianBusiness {
	return func(ctx context.Context, business *Business, clinicianID int) error {
		return edit.EditClinicianBusiness(ctx, business, clinicianID)
	}
}
