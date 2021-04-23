package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func handleGetPatientBookings(getter usecase.PatientBookingsGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handlePatchPreRegisteredBooking(register usecase.BookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handlePostPreRegisteredBooking(register usecase.BookingPreRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handlePostBooking(register usecase.BookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

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

func handlePublicPostBooking(register usecase.BookingRegister) echo.HandlerFunc {
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

func handleDeleteBooking(deleter usecase.BookingSlotDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handleGetFreeBookingSlots(getter usecase.CalendarReader) echo.HandlerFunc {
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

func handleGetBookingSlots(getter usecase.CalendarReader) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handlePostBlockedBookingSlot(blocker usecase.BookingSlotBlocker) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handleDeleteBookingSlotBlocked(deleter usecase.BookingSlotDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
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

func handleDeletePublicBooking(deleter usecase.BookingSlotDeleter) echo.HandlerFunc {
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
