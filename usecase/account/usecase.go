package account

type Usecase struct {
	AccountCreater ClinicianAccountCreater
	AccountGetter  ClinicianAccountGetter

	HomeAddressCreater       ClinicianAddressCreater
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
