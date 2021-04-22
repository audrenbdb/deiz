package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type (
	officeHoursAdder interface {
		AddOfficeHours(ctx context.Context, h *deiz.OfficeHours, clinicianID int) error
	}
	officeHoursRemover interface {
		RemoveOfficeHours(ctx context.Context, hoursID int, clinicianID int) error
	}
)

func handlePostOfficeHours(adder officeHoursAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var h deiz.OfficeHours
		if err := c.Bind(&h); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := adder.AddOfficeHours(ctx, &h, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, h)
	}
}

func handleDeleteOfficeHours(remover officeHoursRemover) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		hoursID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = remover.RemoveOfficeHours(ctx, hoursID, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
