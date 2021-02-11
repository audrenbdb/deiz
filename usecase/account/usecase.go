package account

type accountRepo interface {
	ClinicianAccountCreater
	ClinicianAccountGetter
	ClinicianRegistrationCompleteVerifier
	ClinicianRegistrationCompleter

	ClinicianBusinessUpdater

	ClinicianHomeAddressSetter
	ClinicianHomeAddressCreater
	ClinicianOfficeAddressCreater
	ClinicianAddressOwnershipVerifier
	AddressUpdater

	ClinicianStripeKeysUpdater

	ClinicianGetterByEmail
}

type cryptService interface {
	StringToBytesCrypter
}

type Usecase struct {
	AccountCreater ClinicianAccountCreater
	AccountGetter  ClinicianAccountGetter

	HomeAddressCreater       ClinicianHomeAddressCreater
	HomeAddressSetter        ClinicianHomeAddressSetter
	OfficeAddressCreater     ClinicianOfficeAddressCreater
	AddressOwnerShipVerifier ClinicianAddressOwnershipVerifier
	AddressUpdater           AddressUpdater

	StripeKeysUpdater ClinicianStripeKeysUpdater

	ClinicianGetterByEmail ClinicianGetterByEmail

	RegistrationVerifier  ClinicianRegistrationCompleteVerifier
	RegistrationCompleter ClinicianRegistrationCompleter

	BusinessUpdater ClinicianBusinessUpdater

	StringToBytesCrypter StringToBytesCrypter
}

func NewUsecase(repo accountRepo, cryptsrv cryptService) *Usecase {
	return &Usecase{
		AccountCreater:        repo,
		AccountGetter:         repo,
		RegistrationVerifier:  repo,
		RegistrationCompleter: repo,

		BusinessUpdater: repo,

		HomeAddressSetter:        repo,
		HomeAddressCreater:       repo,
		OfficeAddressCreater:     repo,
		AddressOwnerShipVerifier: repo,
		AddressUpdater:           repo,

		StringToBytesCrypter: cryptsrv,

		ClinicianGetterByEmail: repo,
	}
}
