package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type (
	bookingRegister interface {
		RegisterBookingFromClinician(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error
		RegisterBookingFromPatient(ctx context.Context, b *deiz.Booking) error
		RegisterPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient bool) error
	}
	bookingSlotDeleter interface {
		DeleteBlockedSlot(ctx context.Context, bookingID, clinicianID int) error
		DeletePreRegisteredSlot(ctx context.Context, bookingID, clinicianID int) error
		DeleteBookedSlotFromPatient(ctx context.Context, deleteID string) error
		DeleteBookedSlotFromClinician(ctx context.Context, bookingID int, notifyPatient bool, clinicianID int) error
	}
	bookingPreRegister interface {
		PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error
	}
	bookingSlotBlocker interface {
		BlockBookingSlot(ctx context.Context, slot *deiz.Booking, clinicianID int) error
	}
	calendarReader interface {
		GetCalendarFreeSlots(ctx context.Context, start time.Time, motive deiz.BookingMotive, clinicianID int) ([]deiz.Booking, error)
		GetCalendarSlots(ctx context.Context, start time.Time, motive deiz.BookingMotive, clinicianID int) ([]deiz.Booking, error)
	}

	PatientBookingsGetter interface {
		GetPatientBookings(ctx context.Context, clinicianID, patientID int) ([]deiz.Booking, error)
	}
)

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

func handlePatchPreRegisteredBooking(register bookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		notifyPatient, err := strconv.ParseBool(c.QueryParam("notifyPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var b deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = register.RegisterPreRegisteredBooking(ctx, &b, clinicianID, notifyPatient)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handlePostPreRegisteredBooking(register bookingPreRegister) echo.HandlerFunc {
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

func handlePostBooking(register bookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		notifyPatient, err := strconv.ParseBool(c.QueryParam("notifyPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var b deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = register.RegisterBookingFromClinician(ctx, &b, clinicianID, notifyPatient)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handlePublicPostBooking(register bookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var b deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := register.RegisterBookingFromPatient(ctx, &b); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleDeleteBooking(deleter bookingSlotDeleter) echo.HandlerFunc {
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
		if err := deleter.DeleteBookedSlotFromClinician(ctx, bookingID, notifyPatient, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleGetFreeBookingSlots(getter calendarReader) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID, err := strconv.Atoi(c.QueryParam("clinician"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		from, err := getTimeFromParam(c, "from")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		motive, err := getMotiveFromParam(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		bookings, err := getter.GetCalendarFreeSlots(ctx, from, motive, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handleGetBookingSlots(getter calendarReader) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		from, err := getTimeFromParam(c, "from")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		motive, err := getMotiveFromParam(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		bookings, err := getter.GetCalendarSlots(ctx, from, motive, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func getMotiveFromParam(c echo.Context) (deiz.BookingMotive, error) {
	motiveID, err := strconv.Atoi(c.QueryParam("motiveId"))
	if err != nil {
		return deiz.BookingMotive{}, err
	}
	motiveDuration, err := strconv.Atoi(c.QueryParam("motiveDuration"))
	if err != nil {
		return deiz.BookingMotive{}, err
	}
	return deiz.BookingMotive{ID: motiveID, Duration: motiveDuration}, nil
}

func handlePostBlockedBookingSlot(blocker bookingSlotBlocker) echo.HandlerFunc {
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

func handleDeleteBookingSlotBlocked(deleter bookingSlotDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		bookingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := deleter.DeleteBlockedSlot(ctx, bookingID, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleDeletePublicBooking(deleter bookingSlotDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		deleteID := c.Param("id")
		if len(deleteID) < 6 {
			return c.JSON(http.StatusBadRequest, deiz.ErrorStructValidation)
		}
		err := deleter.DeleteBookedSlotFromPatient(ctx, deleteID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
