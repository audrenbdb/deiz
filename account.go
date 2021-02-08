package deiz

import (
	"context"
	"strings"
)

type ClinicianAccount struct {
	Clinician        Clinician        `json:"clinician"`
	Business         Business         `json:"business"`
	OfficeAddresses  []Address        `json:"officeAddresses"`
	StripePublicKey  string           `json:"stripePublicKey"`
	OfficeHours      []OfficeHours    `json:"officeHours"`
	BookingMotives   []BookingMotive  `json:"bookingMotives"`
	CalendarSettings CalendarSettings `json:"calendarSettings"`
}

type ClinicianAccountRepo struct {
	Getter ClinicianAccountGetter
}

type (
	ClinicianAccountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (ClinicianAccount, error)
	}
)

func (r *Repo) GetClinicianAccount(ctx context.Context, clinicianID int) (ClinicianAccount, error) {
	return ClinicianAccount{}, nil
}

//repo functions
type (
	clinicianAccountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (ClinicianAccount, error)
	}
	clinicianAccountAdder interface {
		AddClinicianAccount(ctx context.Context, c *Clinician) error
	}
	clinicianStripeKeysEditer interface {
		EditClinicianStripeKeys(ctx context.Context, pk string, sk []byte, clinicianID int) error
	}
	stringEncrypter interface {
		EncryptString(ctx context.Context, s string) ([]byte, error)
	}
)

//core functions
type (
	//GetClinicianAccount retrieves clinician account public informations
	GetClinicianAccount func(ctx context.Context, clinicianID int) (ClinicianAccount, error)
	//AddClinicianAccount creates a new clinician user and his default settings
	AddClinicianAccount func(ctx context.Context, c *Clinician) error
	//EditClinicianStripeKeys updates clinician stripe keys
	EditClinicianStripeKeys func(ctx context.Context, pk, sk string, clinicianID int) error
)

func addClinicianAccountFunc(adder clinicianAccountAdder) AddClinicianAccount {
	return func(ctx context.Context, c *Clinician) error {
		return adder.AddClinicianAccount(ctx, &Clinician{
			Name:    strings.ToUpper(c.Name),
			Surname: strings.Title(c.Surname),
			Email:   strings.ToLower(c.Email),
			Phone:   c.Phone,
		})
	}
}

func getClinicianAccountFunc(getter clinicianAccountGetter) GetClinicianAccount {
	return func(ctx context.Context, clinicianID int) (ClinicianAccount, error) {
		return getter.GetClinicianAccount(ctx, clinicianID)
	}
}

func editClinicianStripeKeysFunc(updater clinicianStripeKeysEditer, crypt stringEncrypter) EditClinicianStripeKeys {
	return func(ctx context.Context, pk, sk string, clinicianID int) error {
		skCrypt, err := crypt.EncryptString(ctx, sk)
		if err != nil {
			return err
		}
		return updater.EditClinicianStripeKeys(ctx, pk, skCrypt, clinicianID)
	}
}
