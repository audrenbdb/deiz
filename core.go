package deiz

type crypt interface {
	bytesDecrypter
}

type stripe interface {
	stripePaymentSessionCreater
}

//repo is a driven actor called BY the core to manage storage and persistence
type repo interface {
	clinicianEmailEditer
	clinicianRoleUpdater
	clinicianStripeSecretKeyGetter

	patientsSearcher
	patientsCounter
	patientEditer
	patientAddressEditer
	patientRemover
	patientCreater
}

//Core methods exposed to primary actors
//Core methods to be called FROM external package
type Core struct {
	EditClinicianEmail     EditClinicianEmail
	EnableClinicianAccess  EnableClinicianAccess
	DisableClinicianAccess DisableClinicianAccess

	CreateStripePaymentSession CreateStripePaymentSession

	SearchPatients     SearchPatients
	EditPatient        EditPatient
	RemovePatient      RemovePatient
	EditPatientAddress EditPatientAddress
}

//Implements core function with driven actors
func NewCore(repo repo, crypt crypt, stripe stripe) Core {
	return Core{
		EditClinicianEmail:     editClinicianEmailFunc(repo),
		EnableClinicianAccess:  enableClinicianAccessFunc(repo),
		DisableClinicianAccess: disableClinicianAccessFunc(repo),

		CreateStripePaymentSession: creatStripePaymentSessionFunc(repo, crypt, stripe),

		SearchPatients:     searchPatientsFunc(repo),
		EditPatient:        editPatientFunc(repo),
		EditPatientAddress: editPatientAddressFunc(repo),
		RemovePatient:      removePatientFunc(repo),
	}
}
