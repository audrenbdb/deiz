package echo

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"strconv"
	"time"
)

type (
	BookingUsecases struct {
		Register       bookingRegister
		PreRegister    bookingPreRegister
		CalendarReader calendarReader
		SlotDeleter    bookingSlotDeleter
		SlotBlocker    bookingSlotBlocker
	}
	BillingUsecases struct {
		InvoiceCreater       invoiceCreater
		InvoiceCanceler      invoiceCanceler
		InvoiceMailer        invoiceMailer
		InvoicesGetter       invoicesGetter
		StripeSessionCreater stripeSessionCreater
		UnpaidBookingsGetter unpaidBookingsGetter
	}
	AccountUsecases struct {
		AccountAdder      accountAdder
		LoginAllower      loginAllower
		AccountDataGetter accountDataGetter

		AccountAddressUsecases   AccountAddressUsecases
		BusinessUsecases         businessEditer
		ClinicianUsecases        ClinicianUsecases
		MotiveUsecases           MotiveUsecases
		OfficeHoursUsecases      OfficeHoursUsecases
		CalendarSettingsUsecases calendarSettingsEditer
		StripeKeysUsecases       stripeKeysSetter
	}
	AccountAddressUsecases struct {
		OfficeAddressAdder officeAddressAdder
		AddressDeleter     addressDeleter
		HomeAddressSetter  homeAddressSetter
		AddressEditer      addressEditer
	}
	ClinicianUsecases struct {
		PhoneEditer      clinicianPhoneEditer
		EmailEditer      clinicianEmailEditer
		ProfessionEditer clinicianProfessionEditer
		AdeliEditer      clinicianAdeliEditer
	}
	MotiveUsecases struct {
		MotiveAdder   bookingMotiveAdder
		MotiveEditer  bookingMotiveEditer
		MotiveRemover bookingMotiveRemover
	}
	OfficeHoursUsecases struct {
		OfficeHoursAdder   officeHoursAdder
		OfficeHoursRemover officeHoursRemover
	}
	PatientUsecases interface {
		PatientSearcher
		PatientAdder
		PatientEditer
		PatientAddressAdder
		PatientAddressEditer
		PatientBookingsGetter
	}
	ContactService interface {
		ContactFormToClinicianSender
		GetInTouchSender
	}
)

type EchoServerDeps struct {
	CredentialsGetter credentialsGetter
	AccountUsecases   AccountUsecases
	PatientUsecases   PatientUsecases
	ContactService    ContactService
	BookingUsecases   BookingUsecases
	BillingUsecases   BillingUsecases
}

func StartEchoServer(deps EchoServerDeps) error {
	clinicianMW := roleMW(deps.CredentialsGetter, 2)
	//adminMW := roleMW(credentialsGetter, 3)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	e.POST("/api/registrations", handlePostRegistration(deps.AccountUsecases.LoginAllower))

	e.POST("/api/clinician-accounts", handlePostClinicianAccount(deps.AccountUsecases.AccountAdder))
	e.GET("/api/clinician-accounts/current", handleGetClinicianAccount(deps.AccountUsecases.AccountDataGetter), clinicianMW)
	e.PATCH("/api/businesses/:id", handlePatchBusiness(deps.AccountUsecases.BusinessUsecases), clinicianMW)

	e.GET("/api/bookings", handleGetBookingSlots(deps.BookingUsecases.CalendarReader), clinicianMW)
	e.POST("/api/bookings/blocked", handlePostBlockedBookingSlot(deps.BookingUsecases.SlotBlocker), clinicianMW)
	e.POST("/api/bookings", handlePostBooking(deps.BookingUsecases.Register), clinicianMW)
	e.POST("/api/bookings/pre-registered", handlePostPreRegisteredBooking(deps.BookingUsecases.PreRegister), clinicianMW)
	e.PATCH("/api/bookings/pre-registered", handlePatchPreRegisteredBooking(deps.BookingUsecases.Register), clinicianMW)
	e.DELETE("/api/bookings/:id/blocked", handleDeleteBookingSlotBlocked(deps.BookingUsecases.SlotDeleter), clinicianMW)
	e.DELETE("/api/bookings/:id", handleDeleteBooking(deps.BookingUsecases.SlotDeleter), clinicianMW)

	e.GET("/api/bookings/unpaid", handleGetUnpaidBookings(deps.BillingUsecases.UnpaidBookingsGetter), clinicianMW)

	e.PATCH("/api/clinicians/:id/phone", handlePatchClinicianPhone(deps.AccountUsecases.ClinicianUsecases.PhoneEditer), clinicianMW)
	e.PATCH("/api/clinicians/:id/email", handlePatchClinicianEmail(deps.AccountUsecases.ClinicianUsecases.EmailEditer), clinicianMW)
	e.PATCH("/api/clinicians/:id/adeli", handlePatchClinicianAdeli(deps.AccountUsecases.ClinicianUsecases.AdeliEditer), clinicianMW)
	e.PATCH("/api/clinicians/:id/profession", handlePatchClinicianProfession(deps.AccountUsecases.ClinicianUsecases.ProfessionEditer), clinicianMW)
	e.PATCH("/api/clinicians/:id/addresses/:aid", handlePatchClinicianAddress(deps.AccountUsecases.AccountAddressUsecases.AddressEditer), clinicianMW)
	e.POST("/api/clinicians/:id/addresses", handlePostClinicianAddress(deps.AccountUsecases.AccountAddressUsecases.OfficeAddressAdder, deps.AccountUsecases.AccountAddressUsecases.HomeAddressSetter), clinicianMW)
	e.DELETE("/api/clinicians/:id/addresses/:aid", handleDeleteClinicianAddress(deps.AccountUsecases.AccountAddressUsecases.AddressDeleter), clinicianMW)

	e.GET("/api/patients", handleGetPatients(deps.PatientUsecases), clinicianMW)
	e.POST("/api/patients", handlePostPatient(deps.PatientUsecases), clinicianMW)
	e.PATCH("/api/patients", handlePatchPatient(deps.PatientUsecases), clinicianMW)
	e.POST("/api/patients/:id/address", handlePostPatientAddress(deps.PatientUsecases), clinicianMW)
	e.PATCH("/api/patients/:id/address", handlePatchPatientAddress(deps.PatientUsecases), clinicianMW)
	e.GET("/api/patients/:id/bookings", handleGetPatientBookings(deps.PatientUsecases), clinicianMW)

	e.POST("/api/pdf-booking-invoices/:id", handlePostPDFBookingInvoice(deps.BillingUsecases.InvoiceMailer), clinicianMW)
	e.POST("/api/pdf-booking-invoices", handlePostPDFBookingInvoicesPeriodSummary(deps.BillingUsecases.InvoiceMailer), clinicianMW)
	e.POST("/api/booking-invoices", handlePostBookingInvoice(deps.BillingUsecases.InvoiceCreater), clinicianMW)
	e.POST("/api/booking-invoices/canceled", handlePostCancelInvoice(deps.BillingUsecases.InvoiceCanceler), clinicianMW)
	e.GET("/api/booking-invoices", handleGetPeriodInvoices(deps.BillingUsecases.InvoicesGetter), clinicianMW)

	e.PATCH("/api/clinician-accounts/settings", handlePatchCalendarSettings(deps.AccountUsecases.CalendarSettingsUsecases), clinicianMW)

	e.POST("/api/officehours", handlePostOfficeHours(deps.AccountUsecases.OfficeHoursUsecases.OfficeHoursAdder), clinicianMW)
	e.DELETE("/api/officehours/:id", handleDeleteOfficeHours(deps.AccountUsecases.OfficeHoursUsecases.OfficeHoursRemover), clinicianMW)

	e.POST("/api/booking-motives", handlePostBookingMotive(deps.AccountUsecases.MotiveUsecases.MotiveAdder), clinicianMW)
	e.PATCH("/api/booking-motives/:id", handlePatchBookingMotive(deps.AccountUsecases.MotiveUsecases.MotiveEditer), clinicianMW)
	e.DELETE("/api/booking-motives/:id", handleDeleteBookingMotive(deps.AccountUsecases.MotiveUsecases.MotiveRemover), clinicianMW)

	e.PATCH("/api/clinician-accounts/stripe-keys", handlePatchStripeKeys(deps.AccountUsecases.StripeKeysUsecases), clinicianMW)

	/* PUBLIC API */
	e.GET("/api/public/clinician-accounts", handleGetClinicianAccountPublicData(deps.AccountUsecases.AccountDataGetter))
	e.GET("/api/public/booking-slots", handleGetFreeBookingSlots(deps.BookingUsecases.CalendarReader))
	e.POST("/api/public/bookings", handlePublicPostBooking(deps.BookingUsecases.Register))
	e.GET("/api/public/session-checkout", handleGetSessionCheckout(deps.BillingUsecases.StripeSessionCreater))
	e.DELETE("/api/public/bookings/:id", handleDeletePublicBooking(deps.BookingUsecases.SlotDeleter))
	e.POST("/api/public/contact-form", handlePostContactFormToClinician(deps.ContactService))
	e.POST("/api/public/get-in-touch-form", handlePostGetInTouchForm(deps.ContactService))

	return e.Start(":8080")
}

func FakeCredentialsGetter(ctx context.Context, tokenID string) (credentials, error) {
	return credentials{
		userID: 7,
		role:   2,
	}, nil
}

func getTimeFromParam(c echo.Context, paramName string) (time.Time, error) {
	i, err := strconv.ParseInt(c.QueryParam(paramName), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0).UTC(), nil
}
