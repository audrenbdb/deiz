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
	AccountCreater           ClinicianAccountCreater
	AccountGetter            ClinicianAccountGetter
	ClinicianAddressCreater  ClinicianAddressCreater
	OfficeAddressCreater     ClinicianOfficeAddressCreater
	AddressOwnershipVerifier ClinicianAddressOwnershipVerifier
	AddressUpdater           AddressUpdater
}

type (
	ClinicianAccountCreater interface {
		CreateClinicianAccount(ctx context.Context, account *ClinicianAccount) error
	}
	ClinicianAccountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (ClinicianAccount, error)
	}
	ClinicianOfficeAddressCreater interface {
		CreateClinicianOfficeAddress(ctx context.Context, a *Address, clinicianID int) error
	}
	ClinicianAddressCreater interface {
		CreateClinicianAddress(ctx context.Context, a *Address, clinicianID int) error
	}
	ClinicianAddressOwnershipVerifier interface {
		IsAddressToClinician(ctx context.Context, a *Address, clinicianID int) (bool, error)
	}
)

func (c *Repo) AddClinicianAccount(ctx context.Context, account *ClinicianAccount) error {
	return c.ClinicianAccount.AccountCreater.CreateClinicianAccount(ctx, account)
}

func (c *Repo) GetClinicianAccount(ctx context.Context, clinicianID int) (ClinicianAccount, error) {
	return c.ClinicianAccount.AccountGetter.GetClinicianAccount(ctx, clinicianID)
}

func (c *Repo) AddClinicianOfficeAddress(ctx context.Context, address *Address, clinicianID int) error {
	if !address.isValid() {
		return ErrorStructValidation
	}
	return c.ClinicianAccount.OfficeAddressCreater.CreateClinicianOfficeAddress(ctx, address, clinicianID)
}

func (c *Repo) AddClinicianAddress(ctx context.Context, address *Address, clinicianID int) error {
	if !address.isValid() {
		return ErrorStructValidation
	}
	return c.ClinicianAccount.ClinicianAddressCreater.CreateClinicianAddress(ctx, address, clinicianID)
}

func (c *Repo) UpdateClinicianAddress(ctx context.Context, address *Address, clinicianID int) error {
	if !address.isValid() {
		return ErrorStructValidation
	}
	ownsAddress, err := c.ClinicianAccount.AddressOwnershipVerifier.IsAddressToClinician(ctx, address, clinicianID)
	if err != nil {
		return err
	}
	if !ownsAddress {
		return ErrorUnauthorized
	}
	if err := c.ClinicianAccount.AddressUpdater.UpdateAddress(ctx, address); err != nil {
		return err
	}
	return nil
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
