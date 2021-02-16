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
	BookingSlotsGetter interface {
		GetBookingSlots(ctx context.Context, start time.Time, tzName string, defaultMotiveID, defaultMotiveDuration, clinicianID int) ([]deiz.Booking, error)
	}
	BookingSlotBlocker interface {
		BlockBookingSlot(ctx context.Context, b *deiz.Booking, clinicianID int) error
	}
	BookingSlotUnlocker interface {
		UnlockBookingSlot(ctx context.Context, bookingID, clinicianID int) error
	}
	BookingRegister interface {
		RegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error
	}
	BookingPreRegister interface {
		PreRegisterBooking(ctx context.Context, b *deiz.Booking, clinicianID int) error
	}
	BookingConfirmer interface {
		ConfirmPreRegisteredBooking(ctx context.Context, b *deiz.Booking, clinicianID int, notifyPatient, notifyClinician bool) error
	}
)

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
		from := time.Unix(i, 0)
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

func handleGetBookingsPendingPayment(getBookingsPendingPayments deiz.ListBookingsPendingPayment) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		b, err := getBookingsPendingPayments(ctx, clinicianID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}
