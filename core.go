package deiz

type crypt interface {
	bytesDecrypter
}

type stripe interface {
	stripePaymentSessionCreater
}

type Repo struct {
	Mailing MailingService
}

//repo is a driven actor called BY the core to manage storage and persistence
type repo interface {
	clinicianEmailEditer
	clinicianPhoneEditer
	clinicianRoleUpdater
	clinicianStripeSecretKeyGetter

	bookingMotiveAdder
	bookingMotiveRemover

	patientsSearcher
	patientsCounter
	patientEditer
	patientAddressEditer
	patientRemover
	patientCreater

	officeHoursGetter
	officeHoursAdder
	officeHoursRemover
}

//Core methods exposed to primary actors
//Core methods to be called FROM external package
type Core struct {
	EditClinicianEmail     EditClinicianEmail
	EditClinicianPhone     EditClinicianPhone
	EnableClinicianAccess  EnableClinicianAccess
	DisableClinicianAccess DisableClinicianAccess

	CreateStripePaymentSession CreateStripePaymentSession

	AddBookingMotive    AddBookingMotive
	RemoveBookingMotive RemoveBookingMotive

	SearchPatients     SearchPatients
	EditPatient        EditPatient
	RemovePatient      RemovePatient
	EditPatientAddress EditPatientAddress

	AddOfficeHours    AddOfficeHours
	RemoveOfficeHours RemoveOfficeHours
}

//Implements core function with driven actors
func NewCore(repo repo, crypt crypt, stripe stripe) Core {
	return Core{
		EditClinicianPhone:     editClinicianPhoneFunc(repo),
		EditClinicianEmail:     editClinicianEmailFunc(repo),
		EnableClinicianAccess:  enableClinicianAccessFunc(repo),
		DisableClinicianAccess: disableClinicianAccessFunc(repo),

		CreateStripePaymentSession: creatStripePaymentSessionFunc(repo, crypt, stripe),

		AddBookingMotive:    addBookingMotiveFunc(repo),
		RemoveBookingMotive: removeBookingMotiveFunc(repo),

		SearchPatients:     searchPatientsFunc(repo),
		EditPatient:        editPatientFunc(repo),
		EditPatientAddress: editPatientAddressFunc(repo),
		RemovePatient:      removePatientFunc(repo),

		AddOfficeHours:    addOfficeHoursFunc(repo),
		RemoveOfficeHours: removeOfficeHoursFunc(repo),
	}
}
