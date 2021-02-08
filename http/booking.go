package http

import (
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

func handlePostBookingBlocked(post deiz.FillFreeBookingSlot) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var b *deiz.Booking
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if !b.IsValid() {
			return c.JSON(http.StatusBadRequest, errValidating.Error())
		}
		err := post(ctx, b, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handleDeleteBookingBlocked(remove deiz.FreeBookingSlot) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		bookingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := remove(ctx, &deiz.Booking{ID: bookingID}, clinicianID); err != nil {
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

func handleGetAllBookingsSlot(getSlots deiz.GetAllBookingSlotsFromWeek) echo.HandlerFunc {
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
		bookings, err := getSlots(ctx, from, tzName, motiveID, motiveDuration, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bookings)
	}
}
