package deiz

//pdf is a driven actor called BY the core to create PDF
type pdf interface {
	bookingInvoicePDFGenerater
	periodBookingInvoicesSummaryPDFGetter
}

type crypt interface {
	bytesDecrypter
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

type Repo struct {
	Booking        BookingRepo
	Mailing        MailingService
	GoogleCalendar GoogleCalendarService
	GoogleMaps     GoogleMapsService
	Crypt          CryptService
}

//repo is a driven actor called BY the core to manage storage and persistence
type repo interface {
	clinicianEmailEditer
	clinicianPhoneEditer
	clinicianRoleUpdater
	clinicianInvoicesCounter
	clinicianStripeSecretKeyGetter

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
	patientCreater

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
	EditClinicianEmail     EditClinicianEmail
	EditClinicianPhone     EditClinicianPhone
	EnableClinicianAccess  EnableClinicianAccess
	DisableClinicianAccess DisableClinicianAccess

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
		EditClinicianPhone:     editClinicianPhoneFunc(repo),
		EditClinicianEmail:     editClinicianEmailFunc(repo),
		EnableClinicianAccess:  enableClinicianAccessFunc(repo),
		DisableClinicianAccess: disableClinicianAccessFunc(repo),

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

		FillFreeBookingSlot:         fillFreeBookingSlotFunc(repo, repo),
		FreeBookingSlot:             freeBookingSlotFunc(repo),
		GetAllBookingSlotsFromWeek:  getAllBookingSlotsFromWeekFunc(repo, repo),
		GetFreeBookingSlotsFromWeek: getFreeBookingSlotsFromWeekFunc(repo, repo, repo),
		MailBooking:                 mailBookingFunc(mail, repo),
		MailCancelBooking:           mailCancelBookingFunc(mail, repo),

		AddOfficeHours:    addOfficeHoursFunc(repo),
		RemoveOfficeHours: removeOfficeHoursFunc(repo),
	}
}
