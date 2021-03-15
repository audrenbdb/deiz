package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	AccountService interface {
		ClinicianAccountAdder
		ClinicianAccountGetter
		ClinicianAccountPublicDataGetter
		CalendarSettingsEditer
		ClinicianPhoneEditer
		ClinicianEmailEditer
		ClinicianAddressEditer
		ClinicianAdeliEditer
		ClinicianProfessionEditer
		TaxExemptionCodesGetter
		BusinessEditer
		ClinicianOfficeAddressAdder
		ClinicianHomeAddressAdder
		ClinicianAddressRemover
		OfficeHoursAdder
		OfficeHoursRemover
		BookingMotiveAdder
		BookingMotiveRemover
		BookingMotiveEditer
		StripeKeysSetter
		ClinicianRegistrationCompleter
	}
	BookingService interface {
		BookingSlotsGetter
		FreeBookingSlotsGetter
		BookingSlotBlocker
		BookingSlotUnlocker
		BookingRegister
		PublicBookingRegister
		BookingPreRegister
		BookingConfirmer
		BookingRemover
		PublicBookingRemover
	}
	PatientService interface {
		PatientSearcher
		PatientAdder
		PatientEditer
		PatientAddressAdder
		PatientAddressEditer
		PatientBookingsGetter
	}
	BillingService interface {
		UnpaidBookingsGetter
		BookingInvoiceGenerater
		BookingInvoiceMailer
		PaymentMethodsGetter
		PeriodInvoicesGetter
		PeriodInvoicesSummaryMailer
		StripePaymentSessionCreater
	}
	ContactService interface {
		ContactFormToClinicianSender
		GetInTouchSender
	}
)

func StartEchoServer(
	credentialsGetter credentialsGetter,
	accountService AccountService,
	bookingService BookingService,
	patientService PatientService,
	billingService BillingService,
	contactService ContactService,
) error {
	clinicianMW := roleMW(credentialsGetter, 2)
	//adminMW := roleMW(credentialsGetter, 3)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	e.POST("/api/registrations", handlePostRegistration(accountService))

	e.POST("/api/clinician-accounts", handlePostClinicianAccount(accountService))
	e.GET("/api/clinician-accounts/current", handleGetClinicianAccount(accountService), clinicianMW)
	e.PATCH("/api/businesses/:id", handlePatchBusiness(accountService), clinicianMW)

	e.GET("/api/bookings", handleGetBookingSlots(bookingService), clinicianMW)
	e.POST("/api/bookings/blocked", handlePostBlockedBookingSlot(bookingService), clinicianMW)
	e.POST("/api/bookings", handlePostBooking(bookingService), clinicianMW)
	e.POST("/api/bookings/pre-registered", handlePostPreRegisteredBooking(bookingService), clinicianMW)
	e.PATCH("/api/bookings/pre-registered", handlePatchPreRegisteredBooking(bookingService), clinicianMW)
	e.DELETE("/api/bookings/:id/blocked", handleDeleteBookingSlotBlocked(bookingService), clinicianMW)
	e.DELETE("/api/bookings/:id", handleDeleteBooking(bookingService), clinicianMW)

	e.GET("/api/bookings/unpaid", handleGetUnpaidBookings(billingService), clinicianMW)

	e.PATCH("/api/clinicians/:id/phone", handlePatchClinicianPhone(accountService), clinicianMW)
	e.PATCH("/api/clinicians/:id/email", handlePatchClinicianEmail(accountService), clinicianMW)
	e.PATCH("/api/clinicians/:id/adeli", handlePatchClinicianAdeli(accountService), clinicianMW)
	e.PATCH("/api/clinicians/:id/profession", handlePatchClinicianProfession(accountService), clinicianMW)
	e.PATCH("/api/clinicians/:id/addresses/:aid", handlePatchClinicianAddress(accountService), clinicianMW)
	e.POST("/api/clinicians/:id/addresses", handlePostClinicianAddress(accountService, accountService), clinicianMW)
	e.DELETE("/api/clinicians/:id/addresses/:aid", handleDeleteClinicianAddress(accountService), clinicianMW)

	e.GET("/api/patients", handleGetPatients(patientService), clinicianMW)
	e.POST("/api/patients", handlePostPatient(patientService), clinicianMW)
	e.PATCH("/api/patients", handlePatchPatient(patientService), clinicianMW)
	e.POST("/api/patients/:id/address", handlePostPatientAddress(patientService), clinicianMW)
	e.PATCH("/api/patients/:id/address", handlePatchPatientAddress(patientService), clinicianMW)
	e.GET("/api/patients/:id/bookings", handleGetPatientBookings(patientService), clinicianMW)

	e.POST("/api/pdf-booking-invoices/:id", handlePostPDFBookingInvoice(billingService), clinicianMW)
	e.POST("/api/pdf-booking-invoices", handlePostPDFBookingInvoicesPeriodSummary(billingService), clinicianMW)
	e.POST("/api/booking-invoices", handlePostBookingInvoice(billingService), clinicianMW)
	e.GET("/api/booking-invoices", handleGetPeriodInvoices(billingService), clinicianMW)

	e.PATCH("/api/clinician-accounts/calendar-settings", handlePatchCalendarSettings(accountService), clinicianMW)

	e.GET("/api/payment-methods", handleGetPaymentMethods(billingService), clinicianMW)
	e.GET("/api/tax-exemption-codes", handleGetTaxExemptionCodes(accountService), clinicianMW)

	e.POST("/api/office-hours", handlePostOfficeHours(accountService), clinicianMW)
	e.DELETE("/api/office-hours/:id", handleDeleteOfficeHours(accountService), clinicianMW)

	e.POST("/api/booking-motives", handlePostBookingMotive(accountService), clinicianMW)
	e.PATCH("/api/booking-motives/:id", handlePatchBookingMotive(accountService), clinicianMW)
	e.DELETE("/api/booking-motives/:id", handleDeleteBookingMotive(accountService), clinicianMW)

	e.PATCH("/api/clinician-accounts/stripe-keys", handlePatchStripeKeys(accountService), clinicianMW)

	/* PUBLIC API */
	e.GET("/api/public/clinician-accounts", handleGetClinicianAccountPublicData(accountService))
	e.GET("/api/public/booking-slots", handleGetFreeBookingSlots(bookingService))
	e.POST("/api/public/bookings", handlePublicPostBooking(bookingService))
	e.GET("/api/public/session-checkout", handleGetSessionCheckout(billingService))
	e.DELETE("/api/public/bookings/:id", handleDeletePublicBooking(bookingService))
	e.POST("/api/public/contact-form", handlePostContactFormToClinician(contactService))
	e.POST("/api/public/get-in-touch-form", handlePostGetInTouchForm(contactService))

	return e.Start(":8080")
}

func FakeCredentialsGetter(ctx context.Context, tokenID string) (credentials, error) {
	return credentials{
		userID: 7,
		role:   2,
	}, nil
}
