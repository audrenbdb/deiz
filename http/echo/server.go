package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/auth"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	CredentialsGetter auth.GetCredentialsFromHttpRequest
}

func StartEchoServer(deps EchoServerDeps) error {
	clinicianMW := roleMW(deps.CredentialsGetter, deiz.ClinicianRole)
	publicMW := roleMW(deps.CredentialsGetter, deiz.PublicRole)
	//adminMW := roleMW(credentialsGetter, 3)

	e := echo.New()

	e.POST("/api/registrations", handlePostRegistration(deps.AccountUsecases.LoginAllower))

	e.GET("/api/clinician-accounts", handleGetClinicianAccount(deps.AccountUsecases.AccountDataGetter), publicMW)
	e.POST("/api/clinician-accounts", handlePostClinicianAccount(deps.AccountUsecases.AccountAdder))

	e.PATCH("/api/businesses/:id", handlePatchBusiness(deps.AccountUsecases.BusinessUsecases.BusinessEditer), clinicianMW)
	e.POST("/api/businesses/:bid/address", handlePostBusinessAddress(deps.AccountUsecases.BusinessUsecases.BusinessAddressSetter), clinicianMW)
	e.PATCH("/api/businesses/:bid/addresses/:aid", handlePatchBusinessAddress(deps.AccountUsecases.BusinessUsecases.BusinessAddressEditer), clinicianMW)

	e.GET("/api/bookings", handleGetBookingSlots(deps.BookingUsecases.CalendarReader), clinicianMW)
	e.POST("/api/bookings/blocked", handlePostBlockedBookingSlots(deps.BookingUsecases.SlotBlocker), clinicianMW)
	e.POST("/api/bookings", handlePostBookings(deps.BookingUsecases.Register), clinicianMW)
	e.POST("/api/bookings/pre-registered", handlePostPreRegisteredBookings(deps.BookingUsecases.PreRegister), clinicianMW)
	e.PATCH("/api/bookings/pre-registered", handlePatchPreRegisteredBooking(deps.BookingUsecases.Register), clinicianMW)
	e.DELETE("/api/bookings/:id/blocked", handleDeleteBookingSlotBlocked(deps.BookingUsecases.SlotDeleter), clinicianMW)
	e.DELETE("/api/bookings/:id", handleDeleteBooking(deps.BookingUsecases.SlotDeleter), clinicianMW)

	e.GET("/api/bookings/unpaid", handleGetUnpaidBookings(deps.BillingUsecases.UnpaidBookingsGetter), clinicianMW)

	e.PATCH("/api/clinicians/:id/phone", handlePatchClinicianPhone(deps.AccountUsecases.ClinicianUsecases.PhoneEditer), clinicianMW)
	e.PATCH("/api/clinicians/:id/email", handlePatchClinicianEmail(deps.AccountUsecases.ClinicianUsecases.EmailEditer), clinicianMW)
	e.PATCH("/api/clinicians/:id/adeli", handlePatchClinicianAdeli(deps.AccountUsecases.ClinicianUsecases.AdeliEditer), clinicianMW)
	e.PATCH("/api/clinicians/:id/profession", handlePatchClinicianProfession(deps.AccountUsecases.ClinicianUsecases.ProfessionEditer), clinicianMW)
	e.DELETE("/api/clinicians/:id/addresses/:aid", handleDeleteClinicianAddress(deps.AccountUsecases.AccountAddressUsecases.AddressDeleter), clinicianMW)

	e.POST("/api/office-addresses", handlePostClinicianAddress(deps.AccountUsecases.AccountAddressUsecases.OfficeAddressAdder), clinicianMW)

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

	/* PublicRole API */
	e.GET("/api/public/clinician-accounts", handleGetClinicianAccount(deps.AccountUsecases.AccountDataGetter), publicMW)
	e.GET("/api/public/booking-slots", handleGetFreeBookingSlots(deps.BookingUsecases.CalendarReader))
	e.POST("/api/public/bookings", handlePublicPostBooking(deps.BookingUsecases.Register))
	e.GET("/api/public/session-checkout", handleGetSessionCheckout(deps.BillingUsecases.StripeSessionCreater))
	e.DELETE("/api/public/bookings/:id", handleDeletePublicBooking(deps.BookingUsecases.SlotDeleter))
	e.POST("/api/public/contact-form", handlePostContactFormToClinician(deps.ContactService))
	e.POST("/api/public/get-in-touch-form", handlePostGetInTouchForm(deps.ContactService))

	e.Use(
		middleware.CORS(),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)),
		middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}),
	)
	return e.Start(":8080")
}

func getTimeFromParam(c echo.Context, paramName string) (time.Time, error) {
	i, err := strconv.ParseInt(c.QueryParam(paramName), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0).UTC(), nil
}

func getURLIntegerQueryParam(c echo.Context, paramName string) (int, error) {
	return strconv.Atoi(c.QueryParam(paramName))
}

func getURLIntegerParam(c echo.Context, paramName string) (int, error) {
	return strconv.Atoi(c.Param(paramName))
}
