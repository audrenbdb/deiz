package http

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type (
	FreeBookingSlotsGetter interface {
		GetFreeBookingSlots(ctx context.Context, start time.Time, tzName string, defaultMotiveID, defaultMotiveDuration, clinicianID int) ([]deiz.Booking, error)
	}
	BookingSlotsGetter interface {
		GetBookingSlots(ctx context.Context, start time.Time, tzName string, defaultMotiveID, defaultMotiveDuration, clinicianID int) ([]deiz.Booking, error)
	}
	BookingSlotBlocker interface {
		BlockBookingSlot(ctx context.Context, b *deiz.Booking, clinicianID int) error
	}
	BookingSlotUnlocker interface {
		UnlockBookingSlot(ctx context.Context, bookingID, clinicianID int) error
	}
	PublicBookingRegister interface {
		RegisterPublicBooking(ctx context.Context, b *deiz.Booking) error
	}
	BookingRegister interface {
		RegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error
	}
	BookingRemover interface {
		RemoveBooking(ctx context.Context, bookingID int, notifyPatient bool, clinicianID int) error
	}
	PublicBookingRemover interface {
		RemovePublicBooking(ctx context.Context, deleteID string) error
	}
	BookingPreRegister interface {
		PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error
	}
	BookingConfirmer interface {
		ConfirmPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error
	}
	PatientBookingsGetter interface {
		GetPatientBookings(ctx context.Context, clinicianID, patientID int) ([]deiz.Booking, error)
	}
	UnpaidBookingsGetter interface {
		GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error)
	}
)

func handleGetUnpaidBookings(getter UnpaidBookingsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		bookings, err := getter.GetUnpaidBookings(ctx, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handleGetPatientBookings(getter PatientBookingsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		patientID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		bookings, err := getter.GetPatientBookings(ctx, clinicianID, patientID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handlePatchPreRegisteredBooking(confirmer BookingConfirmer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		notifyPatient, err := strconv.ParseBool(c.QueryParam("notifyPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		notifyClinician, err := strconv.ParseBool(c.QueryParam("notifyClinician"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var b deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = confirmer.ConfirmPreRegisteredBooking(ctx, &b, clinicianID, notifyPatient, notifyClinician)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handlePostPreRegisteredBooking(register BookingPreRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var b deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := register.PreRegisterBooking(ctx, &b, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handlePostBooking(register BookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		notifyPatient, err := strconv.ParseBool(c.QueryParam("notifyPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		notifyClinician, err := strconv.ParseBool(c.QueryParam("notifyClinician"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var b deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = register.RegisterBooking(ctx, &b, clinicianID, notifyPatient, notifyClinician)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handlePublicPostBooking(register PublicBookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var b deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := register.RegisterPublicBooking(ctx, &b); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleDeleteBooking(remover BookingRemover) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		bookingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		notifyPatient, err := strconv.ParseBool(c.QueryParam("notifyPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := remover.RemoveBooking(ctx, bookingID, notifyPatient, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleGetFreeBookingSlots(getter FreeBookingSlotsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID, err := strconv.Atoi(c.QueryParam("clinician"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		i, err := strconv.ParseInt(c.QueryParam("from"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		tzName := c.QueryParam("tz")
		motiveID, err := strconv.Atoi(c.QueryParam("motiveId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		motiveDuration, err := strconv.Atoi(c.QueryParam("motiveDuration"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		from := time.Unix(i, 0).UTC()
		bookings, err := getter.GetFreeBookingSlots(ctx, from, tzName, motiveID, motiveDuration, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handleGetBookingSlots(getter BookingSlotsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		i, err := strconv.ParseInt(c.QueryParam("from"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		tzName := c.QueryParam("tz")
		motiveID, err := strconv.Atoi(c.QueryParam("motiveId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		motiveDuration, err := strconv.Atoi(c.QueryParam("motiveDuration"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		from := time.Unix(i, 0).UTC()
		bookings, err := getter.GetBookingSlots(ctx, from, tzName, motiveID, motiveDuration, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handlePostBlockedBookingSlot(blocker BookingSlotBlocker) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var b *deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := blocker.BlockBookingSlot(ctx, b, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handleDeleteBookingSlotBlocked(unlocker BookingSlotUnlocker) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		bookingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := unlocker.UnlockBookingSlot(ctx, bookingID, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleDeletePublicBooking(remover PublicBookingRemover) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		deleteID := c.Param("id")
		if len(deleteID) < 6 {
			return c.JSON(http.StatusBadRequest, deiz.ErrorStructValidation)
		}
		err := remover.RemovePublicBooking(ctx, deleteID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
