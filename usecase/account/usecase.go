package account

type accountRepo interface {
	ClinicianAccountCreater
	ClinicianAccountGetter
	PublicDataGetter
	ClinicianRegistrationCompleteVerifier
	ClinicianRegistrationCompleter
	ClinicianPhoneUpdater
	ClinicianEmailUpdater
	ClinicianAdeliUpdater
	ClinicianProfessionUpdater

	BookingMotiveUpdater
	BookingMotiveDeleter
	BookingMotiveCreater

	OfficeHoursCreater
	OfficeHoursDeleter

	ClinicianBusinessUpdater
	TaxExemptionCodesGetter
	CalendarSettingsUpdater

	ClinicianHomeAddressSetter
	ClinicianHomeAddressCreater
	ClinicianOfficeAddressCreater
	ClinicianAddressOwnershipVerifier
	AddressUpdater
	AddressDeleter

	ClinicianStripeKeysUpdater

	ClinicianGetterByEmail
}

type cryptService interface {
	StringToBytesCrypter
}

type Usecase struct {
	AccountCreater             ClinicianAccountCreater
	AccountGetter              ClinicianAccountGetter
	PublicDataGetter           PublicDataGetter
	ClinicianPhoneUpdater      ClinicianPhoneUpdater
	ClinicianEmailUpdater      ClinicianEmailUpdater
	ClinicianAdeliUpdater      ClinicianAdeliUpdater
	ClinicianProfessionUpdater ClinicianProfessionUpdater

	OfficeHoursCreater OfficeHoursCreater
	OfficeHoursDeleter OfficeHoursDeleter

	BookingMotiveUpdater BookingMotiveUpdater
	BookingMotiveDeleter BookingMotiveDeleter
	BookingMotiveCreater BookingMotiveCreater

	HomeAddressCreater       ClinicianHomeAddressCreater
	HomeAddressSetter        ClinicianHomeAddressSetter
	OfficeAddressCreater     ClinicianOfficeAddressCreater
	AddressOwnerShipVerifier ClinicianAddressOwnershipVerifier
	AddressUpdater           AddressUpdater
	AddressDeleter           AddressDeleter

	StripeKeysUpdater ClinicianStripeKeysUpdater

	ClinicianGetterByEmail ClinicianGetterByEmail

	RegistrationVerifier  ClinicianRegistrationCompleteVerifier
	RegistrationCompleter ClinicianRegistrationCompleter

	BusinessUpdater         ClinicianBusinessUpdater
	TaxExemptionCodesGetter TaxExemptionCodesGetter

	CalendarSettingsUpdater CalendarSettingsUpdater

	StringToBytesCrypter StringToBytesCrypter
}

func NewUsecase(repo accountRepo, cryptsrv cryptService) *Usecase {
	return &Usecase{
		AccountCreater:             repo,
		AccountGetter:              repo,
		PublicDataGetter:           repo,
		RegistrationVerifier:       repo,
		RegistrationCompleter:      repo,
		ClinicianPhoneUpdater:      repo,
		ClinicianEmailUpdater:      repo,
		ClinicianAdeliUpdater:      repo,
		ClinicianProfessionUpdater: repo,

		OfficeHoursCreater: repo,
		OfficeHoursDeleter: repo,

		BookingMotiveUpdater: repo,
		BookingMotiveDeleter: repo,
		BookingMotiveCreater: repo,

		BusinessUpdater:         repo,
		TaxExemptionCodesGetter: repo,

		HomeAddressSetter:        repo,
		HomeAddressCreater:       repo,
		OfficeAddressCreater:     repo,
		AddressDeleter:           repo,
		AddressOwnerShipVerifier: repo,
		AddressUpdater:           repo,
		CalendarSettingsUpdater:  repo,

		StringToBytesCrypter: cryptsrv,

		ClinicianGetterByEmail: repo,

		StripeKeysUpdater: repo,
	}
}
