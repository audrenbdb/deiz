package deiz

import "context"

type Address struct {
	ID       int    `json:"id" validate:"required"`
	Line     string `json:"line" validate:"required"`
	PostCode int    `json:"postCode" validate:"min=10000"`
	City     string `json:"city" validate:"required"`
}

//repo functions
type (
	clinicianAddressAdder interface {
		AddClinicianPersonalAddress(ctx context.Context, address *Address, clinicianID int) error
		AddClinicianOfficeAddress(ctx context.Context, address *Address, clinicianID int) error
	}
	clinicianAddressEditer interface {
		EditClinicianAddress(ctx context.Context, address *Address, clinicianID int) error
	}
)

//core functions
type (
	//AddClinicianOfficeAddress adds an address where the clinician works
	AddClinicianOfficeAddress func(ctx context.Context, address *Address, clinicianID int) error
	//AddClinicianPersonalAddress adds clinician personal address
	AddClinicianPersonalAddress func(ctx context.Context, address *Address, clinicianID int) error
	//EditClinicianAddress edits a clinician existing address
	EditClinicianAddress func(ctx context.Context, address *Address, clinicianID int) error
)

func addClinicianOfficeAddressFunc(adder clinicianAddressAdder) AddClinicianOfficeAddress {
	return func(ctx context.Context, address *Address, clinicianID int) error {
		return adder.AddClinicianOfficeAddress(ctx, address, clinicianID)
	}
}

func addClinicianPersonalAddressFunc(adder clinicianAddressAdder) AddClinicianPersonalAddress {
	return func(ctx context.Context, address *Address, clinicianID int) error {
		return adder.AddClinicianPersonalAddress(ctx, address, clinicianID)
	}
}

func editClinicianAddressFunc(editer clinicianAddressEditer) EditClinicianAddress {
	return func(ctx context.Context, address *Address, clinicianID int) error {
		return editer.EditClinicianAddress(ctx, address, clinicianID)
	}
}
