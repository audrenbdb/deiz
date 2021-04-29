package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
)

func handlePatchBookingMotive(editer usecase.BookingMotiveEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

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

func handlePostBookingMotive(adder usecase.BookingMotiveAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

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

func handleDeleteBookingMotive(remover usecase.BookingMotiveRemover) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

		motiveID, err := getURLIntegerParam(c, "id")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := remover.RemoveBookingMotive(ctx, motiveID, clinicianID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
