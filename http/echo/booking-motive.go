package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type (
	bookingMotiveEditer interface {
		EditBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
	bookingMotiveRemover interface {
		RemoveBookingMotive(ctx context.Context, mID, clinicianID int) error
	}
	bookingMotiveAdder interface {
		AddBookingMotive(ctx context.Context, m *deiz.BookingMotive, clinicianID int) error
	}
)

func handlePatchBookingMotive(editer bookingMotiveEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		var m deiz.BookingMotive
		if err := c.Bind(&m); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := editer.EditBookingMotive(ctx, &m, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostBookingMotive(adder bookingMotiveAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		var m deiz.BookingMotive
		if err := c.Bind(&m); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := adder.AddBookingMotive(ctx, &m, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, m)
	}
}

func handleDeleteBookingMotive(remover bookingMotiveRemover) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

		motiveID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := remover.RemoveBookingMotive(ctx, motiveID, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
