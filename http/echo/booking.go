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
		patientID, err := getURLIntegerParam(c, "id")
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

func handlePostPreRegisteredBookings(register usecase.BookingPreRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		var bookings []*deiz.Booking
		if err := c.Bind(&bookings); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := register.PreRegisterBookings(ctx, bookings, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}

func handlePostBookings(register usecase.BookingRegister) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

		notifyPatient, err := strconv.ParseBool(c.QueryParam("notifyPatient"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var bookings []*deiz.Booking
		if err := c.Bind(&bookings); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = register.RegisterBookingsFromClinician(ctx, bookings, clinicianID, notifyPatient)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
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
		bookingID, err := getURLIntegerParam(c, "id")
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
		clinicianID, err := getURLIntegerQueryParam(c, "clinician")
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
	motiveID, err := getURLIntegerQueryParam(c, "motiveId")
	if err != nil {
		return deiz.BookingMotive{}, err
	}
	motiveDuration, err := getURLIntegerQueryParam(c, "motiveDuration")
	if err != nil {
		return deiz.BookingMotive{}, err
	}
	return deiz.BookingMotive{ID: motiveID, Duration: motiveDuration}, nil
}

func handlePostBlockedBookingSlots(blocker usecase.BookingSlotBlocker) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cred := getCredFromEchoCtx(c)
		var slots []*deiz.Booking
		if err := c.Bind(&slots); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := blocker.BlockBookingSlots(ctx, slots, cred)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, slots)
	}
}

func handleDeleteBookingSlotBlocked(deleter usecase.BookingSlotDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		bookingID, err := getURLIntegerParam(c, "id")
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
