package http

import (
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartEchoServer(credentialsGetter credentialsGetter, core deiz.Core, v validater) error {
	clinicianMW := roleMW(credentialsGetter, 2)
	adminMW := roleMW(credentialsGetter, 3)

	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	e.POST("/api/credentials", handleLogin(core.Login))

	e.POST("/api/clinicians", handlePostClinician(core.AddClinicianAccount, v), adminMW)

	e.GET("/api/clinicians/:id/account", handleGetClinicianAccount(core.GetClinicianAccount))

	e.POST("/api/clinicians/:id/addresses", handlePostClinicianAddress(core.AddClinicianPersonalAddress, core.AddClinicianOfficeAddress, v), clinicianMW)
	e.PATCH("/api/clinicians/:id/addresses/:aid", handlePatchClinicianAddress(core.EditClinicianAddress, v), clinicianMW)

	e.PATCH("/api/clinicians/:id/calendar-settings/:cid", handlePatchCalendarSettings(core.EditCalendarSettings, v), clinicianMW)

	e.PATCH("/api/clinicians/:id/phone", handlePatchClinicianPhone(core.EditClinicianPhone, v), clinicianMW)
	e.PATCH("/api/clinicians/:id/email", handlePatchClinicianEmail(core.EditClinicianEmail, v), clinicianMW)

	e.GET("/api/clinicians/:id/bookings/pending-payment", handleGetBookingsPendingPayment(core.ListBookingsPendingPayment), clinicianMW)

	e.GET("/api/clinicians/:id/booking-invoices/:bid/pdf", handleGetBookingInvoicePDF(core.SeeInvoicePDF, v), clinicianMW)
	e.POST("/api/clinicians/:id/booking-invoices", handlePostBookingInvoice(core.CreateBookingInvoice, core.MailBookingInvoice, v), clinicianMW)

	e.POST("/api/clinicians/:id/booking-motives", handlePostBookingMotive(core.AddBookingMotive, v), clinicianMW)
	e.DELETE("/api/clinicians/:id/booking-motives/:id", handleDeleteBookingMotive(core.RemoveBookingMotive, v), clinicianMW)

	e.GET("/api/clinicians/:id/patients", handleGetPatients(core.SearchPatients), clinicianMW)
	e.PATCH("/api/clinicians/:id/patients/:pid", handlePatchPatient(core.EditPatient, v), clinicianMW)

	return e.Start("localhost:8080")
}
