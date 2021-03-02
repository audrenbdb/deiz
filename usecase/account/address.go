package account

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	ClinicianOfficeAddressCreater interface {
		CreateClinicianOfficeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
	ClinicianHomeAddressCreater interface {
		CreateClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
	ClinicianHomeAddressSetter interface {
		SetClinicianHomeAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
	ClinicianAddressOwnershipVerifier interface {
		IsAddressToClinician(ctx context.Context, a *deiz.Address, clinicianID int) (bool, error)
	}
	AddressUpdater interface {
		UpdateAddress(ctx context.Context, address *deiz.Address) error
	}
	AddressDeleter interface {
		DeleteAddress(ctx context.Context, addressID int) error
	}
)

func IsAddressValid(a *deiz.Address) bool {
	if len(a.Line) < 2 || a.PostCode < 10000 || len(a.City) < 2 {
		return false
	}
	return true
}

func (u *Usecase) AddClinicianOfficeAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	if !IsAddressValid(address) {
		return deiz.ErrorStructValidation
	}
	return u.OfficeAddressCreater.CreateClinicianOfficeAddress(ctx, address, clinicianID)
}

func (u *Usecase) AddClinicianHomeAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	if !IsAddressValid(address) {
		return deiz.ErrorStructValidation
	}
	if address.ID == 0 {
		return u.HomeAddressCreater.CreateClinicianHomeAddress(ctx, address, clinicianID)
	}
	return u.HomeAddressSetter.SetClinicianHomeAddress(ctx, address, clinicianID)
}

func (u *Usecase) EditClinicianAddress(ctx context.Context, address *deiz.Address, clinicianID int) error {
	if !IsAddressValid(address) {
		return deiz.ErrorStructValidation
	}
	ownsAddress, err := u.AddressOwnerShipVerifier.IsAddressToClinician(ctx, address, clinicianID)
	if err != nil {
		return err
	}
	if !ownsAddress {
		return deiz.ErrorUnauthorized
	}
	if err := u.AddressUpdater.UpdateAddress(ctx, address); err != nil {
		return err
	}
	return nil
}

func (u *Usecase) RemoveClinicianAddress(ctx context.Context, addressID int, clinicianID int) error {
	ownsAddress, err := u.AddressOwnerShipVerifier.IsAddressToClinician(ctx, &deiz.Address{ID: addressID}, clinicianID)
	if err != nil {
		return err
	}
	if !ownsAddress {
		return deiz.ErrorUnauthorized
	}
	return u.AddressDeleter.DeleteAddress(ctx, addressID)
}
