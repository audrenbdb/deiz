package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
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
		credentials := getCredFromEchoCtx(c)
		var acc deiz.ClinicianAccount
		var err error
		if credentials.Role == deiz.ClinicianRole {
			acc, err = getter.GetClinicianAccountData(ctx, credentials)
		} else {
			clinicianID, err := getURLIntegerQueryParam(c, "clinicianId")
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			acc, err = getter.GetClinicianAccountPublicData(ctx, clinicianID)
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, acc)
	}
}
