package deiz

import (
	"context"
	"strings"
)

type Clinician struct {
	ID         int     `json:"id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Surname    string  `json:"surname" validate:"required"`
	Phone      string  `json:"phone" validate:"required,min=10"`
	Email      string  `json:"email" validate:"required,email"`
	Address    Address `json:"personalAddress" validate:"required"`
	Profession string  `json:"profession"`
	Adeli      Adeli   `json:"adeli"`
}

type Adeli struct {
	ID         int    `json:"id" validate:"required"`
	Identifier string `json:"identifier" validate:"required,min=8"`
}

type (
	clinicianRoleUpdater interface {
		UpdateClinicianRole(ctx context.Context, rol, clinicianID int) error
	}
	clinicianEmailEditer interface {
		EditClinicianEmail(ctx context.Context, email string, clinicianID int) error
	}
	clinicianPhoneEditer interface {
		EditClinicianPhone(ctx context.Context, phone string, clinicianID int) error
	}
)

type (
	//EnableClinicianAccess enables clinician access to the application
	EnableClinicianAccess func(ctx context.Context, clinicianID int) error
	//DisableClinicianAccess revokes credentials of a given clinician
	DisableClinicianAccess func(ctx context.Context, clinicianID int) error
	//EditClinicianEmail update email of a given clinician
	EditClinicianEmail func(ctx context.Context, email string, clinicianID int) error
	//EditClinicianPhone update phone of a given clinician
	EditClinicianPhone func(ctx context.Context, phone string, clinicianID int) error
)

func enableClinicianAccessFunc(updater clinicianRoleUpdater) EnableClinicianAccess {
	return func(ctx context.Context, clinicianID int) error {
		return updater.UpdateClinicianRole(ctx, 2, clinicianID)
	}
}

func disableClinicianAccessFunc(updater clinicianRoleUpdater) DisableClinicianAccess {
	return func(ctx context.Context, clinicianID int) error {
		return updater.UpdateClinicianRole(ctx, 1, clinicianID)
	}
}

func editClinicianEmailFunc(editer clinicianEmailEditer) EditClinicianEmail {
	return func(ctx context.Context, email string, clinicianID int) error {
		return editer.EditClinicianEmail(ctx, strings.ToLower(email), clinicianID)
	}
}

func editClinicianPhoneFunc(editer clinicianPhoneEditer) EditClinicianPhone {
	return func(ctx context.Context, phone string, clinicianID int) error {
		return editer.EditClinicianPhone(ctx, phone, clinicianID)
	}
}
