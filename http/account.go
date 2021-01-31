package http

import (
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func handleGetClinicianAccount(getAccount deiz.GetClinicianAccount) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.QueryParam("clinicianId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		acc, err := getAccount(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, acc)
	}
}
