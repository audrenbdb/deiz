package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func handlePostOfficeHours(adder usecase.OfficeHoursAdder) echo.HandlerFunc {
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

func handleDeleteOfficeHours(remover usecase.OfficeHoursRemover) echo.HandlerFunc {
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
