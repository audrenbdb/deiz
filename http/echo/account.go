package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func handlePostRegistration(allower usecase.LoginAllower) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var f deiz.LoginData
		if err := c.Bind(&f); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := allower.AllowLogin(ctx, f)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostClinicianAccount(adder usecase.AccountAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var acc deiz.ClinicianAccount
		if err := c.Bind(&acc); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := adder.AddAccount(ctx, &acc); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleGetClinicianAccount(getter usecase.AccountDataGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID, err := strconv.Atoi(c.QueryParam("clinicianId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		acc, err := getter.GetClinicianAccountData(ctx, clinicianID, getCredFromEchoCtx(c))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, acc)
	}
}
