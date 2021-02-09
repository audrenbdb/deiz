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

//repo functions
type (
	clinicianAddressAdder interface {
		AddClinicianPersonalAddress(ctx context.Context, address *Address, clinicianID int) error
		CreateClinicianOfficeAddress(ctx context.Context, address *Address, clinicianID int) error
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
		return adder.CreateClinicianOfficeAddress(ctx, address, clinicianID)
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
