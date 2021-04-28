package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"strconv"
	"time"
)

type (
	ContactService interface {
		ContactFormToClinicianSender
		GetInTouchSender
	}
)

type EchoServerDeps struct {
	AccountUsecases   usecase.AccountUsecases
	PatientUsecases   usecase.PatientUsecases
	BookingUsecases   usecase.BookingUsecases
	BillingUsecases   usecase.BillingUsecases
	ContactService    ContactService
	CredentialsGetter credentialsGetter
}

func StartEchoServer(deps EchoServerDeps) error {
	clinicianMW := roleMW(deps.CredentialsGetter, 1)
	//adminMW := roleMW(credentialsGetter, 3)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	e.POST("/api/registrations", handlePostRegistration(deps.AccountUsecases.LoginAllower))

	e.POST("/api/clinician-accounts", handlePostClinicianAccount(deps.AccountUsecases.AccountAdder))
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

	e.GET("/api/patients", handleGetPatients(deps.PatientUsecases.Searcher), clinicianMW)
	e.POST("/api/patients", handlePostPatient(deps.PatientUsecases.Adder), clinicianMW)
	e.PATCH("/api/patients", handlePatchPatient(deps.PatientUsecases.Editer), clinicianMW)
	e.POST("/api/patients/:id/address", handlePostPatientAddress(deps.PatientUsecases.AddressAdder), clinicianMW)
	e.PATCH("/api/patients/:id/address", handlePatchPatientAddress(deps.PatientUsecases.AddressEditer), clinicianMW)
	e.GET("/api/patients/:id/bookings", handleGetPatientBookings(deps.PatientUsecases.BookingsGetter), clinicianMW)

	e.POST("/api/pdf-booking-invoices/:id", handlePostPDFBookingInvoice(deps.BillingUsecases.InvoiceMailer), clinicianMW)
	e.POST("/api/pdf-booking-invoices", handlePostPDFBookingInvoicesPeriodSummary(deps.BillingUsecases.InvoiceMailer), clinicianMW)
	e.POST("/api/booking-invoices", handlePostBookingInvoice(deps.BillingUsecases.InvoiceCreater), clinicianMW)
	e.POST("/api/booking-invoices/canceled", handlePostCancelInvoice(deps.BillingUsecases.InvoiceCanceler), clinicianMW)
	e.GET("/api/booking-invoices", handleGetPeriodInvoices(deps.BillingUsecases.InvoicesGetter), clinicianMW)

	e.PATCH("/api/clinician-accounts/calendar-settings", handlePatchCalendarSettings(deps.AccountUsecases.CalendarSettingsUsecases), clinicianMW)

	e.POST("/api/office-hours", handlePostOfficeHours(deps.AccountUsecases.OfficeHoursUsecases.OfficeHoursAdder), clinicianMW)
	e.DELETE("/api/office-hours/:id", handleDeleteOfficeHours(deps.AccountUsecases.OfficeHoursUsecases.OfficeHoursRemover), clinicianMW)

	e.POST("/api/booking-motives", handlePostBookingMotive(deps.AccountUsecases.MotiveUsecases.MotiveAdder), clinicianMW)
	e.PATCH("/api/booking-motives/:id", handlePatchBookingMotive(deps.AccountUsecases.MotiveUsecases.MotiveEditer), clinicianMW)
	e.DELETE("/api/booking-motives/:id", handleDeleteBookingMotive(deps.AccountUsecases.MotiveUsecases.MotiveRemover), clinicianMW)

	e.PATCH("/api/clinician-accounts/stripe-keys", handlePatchStripeKeys(deps.AccountUsecases.StripeKeysUsecases), clinicianMW)

	/* PUBLIC API */
	e.GET("/api/clinician-accounts", handleGetClinicianAccount(deps.AccountUsecases.AccountDataGetter), publicMW(deps.CredentialsGetter))
	e.GET("/api/public/booking-slots", handleGetFreeBookingSlots(deps.BookingUsecases.CalendarReader))
	e.POST("/api/public/bookings", handlePublicPostBooking(deps.BookingUsecases.Register))
	e.GET("/api/public/session-checkout", handleGetSessionCheckout(deps.BillingUsecases.StripeSessionCreater))
	e.DELETE("/api/public/bookings/:id", handleDeletePublicBooking(deps.BookingUsecases.SlotDeleter))
	e.POST("/api/public/contact-form", handlePostContactFormToClinician(deps.ContactService))
	e.POST("/api/public/get-in-touch-form", handlePostGetInTouchForm(deps.ContactService))

	return e.Start(":8080")
}

func FakeCredentialsGetter(ctx context.Context, tokenID string) (deiz.Credentials, error) {
	return deiz.Credentials{
		UserID: 7,
		Role:   deiz.Role(1),
	}, nil
}

func getTimeFromParam(c echo.Context, paramName string) (time.Time, error) {
	i, err := strconv.ParseInt(c.QueryParam(paramName), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0).UTC(), nil
}
