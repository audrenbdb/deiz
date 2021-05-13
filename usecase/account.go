/*
Package usecase references all usecases to be implemented
*/
package usecase

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	AccountUsecases struct {
		AccountAdder      AccountAdder
		LoginAllower      LoginAllower
		AccountDataGetter AccountDataGetter

		AccountAddressUsecases   AccountAddressUsecases
		BusinessUsecases         BusinessUsecases
		ClinicianUsecases        ClinicianUsecases
		MotiveUsecases           MotiveUsecases
		OfficeHoursUsecases      OfficeHoursUsecases
		CalendarSettingsUsecases CalendarSettingsEditer
		StripeKeysUsecases       StripeKeysSetter
	}
	AccountAddressUsecases struct {
		OfficeAddressAdder OfficeAddressAdder
		AddressDeleter     AddressDeleter
		AddressEditer      AddressEditer
	}
	BusinessUsecases struct {
		BusinessEditer        BusinessEditer
		BusinessAddressEditer BusinessAddressEditer
		BusinessAddressSetter BusinessAddressSetter
	}
	ClinicianUsecases struct {
		PhoneEditer      ClinicianPhoneEditer
		EmailEditer      ClinicianEmailEditer
		ProfessionEditer ClinicianProfessionEditer
		AdeliEditer      ClinicianAdeliEditer
	}
	MotiveUsecases struct {
		MotiveAdder   BookingMotiveAdder
		MotiveEditer  BookingMotiveEditer
		MotiveRemover BookingMotiveRemover
	}
	OfficeHoursUsecases struct {
		OfficeHoursAdder   OfficeHoursAdder
		OfficeHoursRemover OfficeHoursRemover
	}
)

type (
	AccountDataGetter interface {
		GetClinicianAccountData(ctx context.Context, cred deiz.Credentials) (deiz.ClinicianAccount, error)
		GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error)
	}
	AccountAdder interface {
		AddAccount(ctx context.Context, acc *deiz.ClinicianAccount) error
	}
	LoginAllower interface {
		AllowLogin(ctx context.Context, loginCredentials deiz.LoginData) error
	}
)

type (
	OfficeAddressAdder interface {
		AddClinicianOfficeAddress(ctx context.Context, address *deiz.Address, cred deiz.Credentials) error
	}
	AddressDeleter interface {
		DeleteAddress(ctx context.Context, addressID int, cred deiz.Credentials) error
	}
)

type (
	BookingMotiveEditer interface {
		EditBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
	BookingMotiveRemover interface {
		RemoveBookingMotive(ctx context.Context, mID, clinicianID int) error
	}
	BookingMotiveAdder interface {
		AddBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
)

type (
	BusinessEditer interface {
		EditClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error
	}
	BusinessAddressSetter interface {
		SetClinicianBusinessAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
	BusinessAddressEditer interface {
		UpdateClinicianBusinessAddress(ctx context.Context, a *deiz.Address, clinicianID int) error
	}
)

type (
	CalendarSettingsEditer interface {
		EditCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error
	}
)

type (
	OfficeHoursAdder interface {
		AddOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error
	}
	OfficeHoursRemover interface {
		RemoveOfficeHours(ctx context.Context, hoursID int, clinicianID int) error
	}
)

type (
	StripeKeysSetter interface {
		SetClinicianStripeKeys(ctx context.Context, pk, sk string, clinicianID int) error
	}
)

type (
	ClinicianPhoneEditer interface {
		EditClinicianPhone(ctx context.Context, phone string, clinicianID int) error
	}
	ClinicianEmailEditer interface {
		EditClinicianEmail(ctx context.Context, email string, clinicianID int) error
	}
	AddressEditer interface {
		EditAddress(ctx context.Context, address *deiz.Address, cred deiz.Credentials) error
	}
	ClinicianAdeliEditer interface {
		EditClinicianAdeli(ctx context.Context, identifier string, clinicianID int) error
	}
	ClinicianProfessionEditer interface {
		EditClinicianProfession(ctx context.Context, profession string, clinicianID int) error
	}
)
