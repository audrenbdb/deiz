package deiz

//pdf is a driven actor called BY the core to create PDF
type pdf interface {
	bookingInvoicePDFGenerater
	periodBookingInvoicesSummaryPDFGetter
}

type crypt interface {
	bytesDecrypter
	stringEncrypter
}

type stripe interface {
	stripePaymentSessionCreater
}

//mail is a driven actor called BY the core to send emails
type mail interface {
	bookingInvoiceMailer
	bookingMailer
	bookingCancelMailer
}

//repo is a driven actor called BY the core to manage storage and persistence
type repo interface {
	logger

	clinicianAccountAdder
	clinicianAccountGetter
	clinicianAddressAdder
	clinicianAddressEditer
	clinicianEmailEditer
	clinicianPhoneEditer
	clinicianRoleUpdater
	clinicianInvoicesCounter
	clinicianBusinessEditer
	clinicianStripeSecretKeyGetter
	clinicianStripeKeysEditer
	clinicianTimezoneGetter

	calendarSettingsEditer
	calendarSettingsGetter

	bookingsPendingPaymentGetter
	bookingInvoiceCreater
	periodBookingInvoicesGetter

	bookingMotiveAdder
	bookingMotiveRemover

	patientsSearcher
	patientsCounter
	patientEditer
	patientAddressEditer
	patientRemover

	freeBookingSlotFiller
	bookingSlotRemover
	bookingsInTimeRangeGetter

	officeHoursGetter
	officeHoursAdder
	officeHoursRemover
}

//Core methods exposed to primary actors
//Core methods to be called FROM external package
type Core struct {
	Login Login

	AddClinicianAccount         AddClinicianAccount
	GetClinicianAccount         GetClinicianAccount
	AddClinicianPersonalAddress AddClinicianPersonalAddress
	AddClinicianOfficeAddress   AddClinicianOfficeAddress
	EditClinicianAddress        EditClinicianAddress
	EditClinicianEmail          EditClinicianEmail
	EditClinicianPhone          EditClinicianPhone
	EditClinicianBusiness       EditClinicianBusiness
	EnableClinicianAccess       EnableClinicianAccess
	DisableClinicianAccess      DisableClinicianAccess
	EditClinicianStripeKeys     EditClinicianStripeKeys

	EditCalendarSettings EditCalendarSettings

	ListBookingsPendingPayment  ListBookingsPendingPayment
	SeeInvoicePDF               SeeInvoicePDF
	MailBookingInvoice          MailBookingInvoice
	CreateBookingInvoice        CreateBookingInvoice
	SeePeriodInvoicesSummaryPDF SeePeriodInvoicesSummaryPDF

	CreateStripePaymentSession CreateStripePaymentSession

	AddBookingMotive    AddBookingMotive
	RemoveBookingMotive RemoveBookingMotive

	SearchPatients     SearchPatients
	EditPatient        EditPatient
	RemovePatient      RemovePatient
	EditPatientAddress EditPatientAddress

	FillFreeBookingSlot         FillFreeBookingSlot
	FreeBookingSlot             FreeBookingSlot
	GetAllBookingSlotsFromWeek  GetAllBookingSlotsFromWeek
	GetFreeBookingSlotsFromWeek GetFreeBookingSlotsFromWeek
	MailBooking                 MailBooking
	MailCancelBooking           MailCancelBooking

	AddOfficeHours    AddOfficeHours
	RemoveOfficeHours RemoveOfficeHours
}

//Implements core function with driven actors
func NewCore(repo repo, pdf pdf, mail mail, crypt crypt, stripe stripe) Core {
	return Core{
		Login: login(repo),

		AddClinicianAccount:         addClinicianAccountFunc(repo),
		GetClinicianAccount:         getClinicianAccountFunc(repo),
		AddClinicianPersonalAddress: addClinicianPersonalAddressFunc(repo),
		AddClinicianOfficeAddress:   addClinicianOfficeAddressFunc(repo),
		EditClinicianAddress:        editClinicianAddressFunc(repo),
		EditClinicianPhone:          editClinicianPhoneFunc(repo),
		EditClinicianEmail:          editClinicianEmailFunc(repo),
		EditClinicianBusiness:       editClinicianBusinessFunc(repo),
		EnableClinicianAccess:       enableClinicianAccessFunc(repo),
		DisableClinicianAccess:      disableClinicianAccessFunc(repo),
		EditClinicianStripeKeys:     editClinicianStripeKeysFunc(repo, crypt),

		EditCalendarSettings: editCalendarSettingsFunc(repo),

		ListBookingsPendingPayment:  listBookingsPendingPaymentFunc(repo),
		SeeInvoicePDF:               seeInvoicePDFFunc(pdf),
		MailBookingInvoice:          mailBookingInvoiceFunc(pdf, mail),
		CreateBookingInvoice:        createBookingInvoiceFunc(repo, repo),
		SeePeriodInvoicesSummaryPDF: seePeriodBookingInvoicesSummaryPDFFunc(repo, pdf, repo),

		CreateStripePaymentSession: creatStripePaymentSessionFunc(repo, crypt, stripe),

		AddBookingMotive:    addBookingMotiveFunc(repo),
		RemoveBookingMotive: removeBookingMotiveFunc(repo),

		SearchPatients:     searchPatientsFunc(repo),
		EditPatient:        editPatientFunc(repo),
		EditPatientAddress: editPatientAddressFunc(repo),
		RemovePatient:      removePatientFunc(repo),

		FillFreeBookingSlot:         fillFreeBookingSlotFunc(repo),
		FreeBookingSlot:             freeBookingSlotFunc(repo),
		GetAllBookingSlotsFromWeek:  getAllBookingSlotsFromWeekFunc(repo, repo, repo),
		GetFreeBookingSlotsFromWeek: getFreeBookingSlotsFromWeekFunc(repo, repo, repo),
		MailBooking:                 mailBookingFunc(mail, repo),
		MailCancelBooking:           mailCancelBookingFunc(mail, repo),

		AddOfficeHours:    addOfficeHoursFunc(repo),
		RemoveOfficeHours: removeOfficeHoursFunc(repo),
	}
}
