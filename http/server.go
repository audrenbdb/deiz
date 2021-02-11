package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	AccountService interface {
		ClinicianAccountGetter
	}
	BookingService interface {
		BookingSlotsGetter
		BookingSlotBlocker
		BookingSlotUnlocker
	}
	PatientService interface {
		PatientSearcher
	}
)

func StartEchoServer(
	credentialsGetter credentialsGetter,
	accountService AccountService,
	bookingService BookingService,
	patientService PatientService,
) error {
	clinicianMW := roleMW(credentialsGetter, 2)
	//adminMW := roleMW(credentialsGetter, 3)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	e.GET("/api/clinician-accounts/current", handleGetClinicianAccount(accountService), clinicianMW)

	e.GET("/api/bookings", handleGetBookingSlots(bookingService), clinicianMW)
	e.POST("/api/bookings/blocked", handlePostBlockedBookingSlot(bookingService), clinicianMW)
	e.DELETE("/api/bookings/:id/blocked", handleDeleteBookingSlotBlocked(bookingService), clinicianMW)

	e.GET("/api/patients", HandleGetPatients(patientService), clinicianMW)

	/*

		e.PATCH("/api/clinicians/:id/calendar-settings/:cid", handlePatchCalendarSettings(core.EditCalendarSettings, v), clinicianMW)

		e.PATCH("/api/clinicians/:id/phone", handlePatchClinicianPhone(core.EditClinicianPhone, v), clinicianMW)
		e.PATCH("/api/clinicians/:id/email", handlePatchClinicianEmail(core.EditClinicianEmail, v), clinicianMW)

		e.GET("/api/bookings/pending-payment", handleGetBookingsPendingPayment(core.ListBookingsPendingPayment), clinicianMW)
		e.GET("/api/bookings", handleGetAllBookingsSlot(core.GetAllBookingSlotsFromWeek), clinicianMW)
		e.POST("/api/bookings/blocked", handlePostBookingBlocked(core.FillFreeBookingSlot), clinicianMW)
		e.DELETE("/api/bookings/:id/blocked", handleDeleteBookingBlocked(core.FreeBookingSlot), clinicianMW)

		e.GET("/api/booking-invoices/:id/pdf", handleGetBookingInvoicePDF(core.SeeInvoicePDF, v), clinicianMW)
		e.POST("/api/booking-invoices", handlePostBookingInvoice(core.CreateBookingInvoice, core.MailBookingInvoice, v), clinicianMW)

		e.POST("/api/clinicians/:id/booking-motives", handlePostBookingMotive(core.AddBookingMotive, v), clinicianMW)
		e.DELETE("/api/clinicians/:id/booking-motives/:id", handleDeleteBookingMotive(core.RemoveBookingMotive, v), clinicianMW)

		e.GET("/api/patients", handleGetPatients(core.SearchPatients), clinicianMW)
		e.PATCH("/api/patients/:pid", handlePatchPatient(core.EditPatient, v), clinicianMW)
	*/

	return e.Start(":8080")
}

func FakeCredentialsGetter(ctx context.Context, tokenID string) (credentials, error) {
	return credentials{
		userID: 1,
		role:   2,
	}, nil
}
