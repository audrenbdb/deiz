package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type (
	accountDataGetter interface {
		GetClinicianAccountData(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error)
		GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccountPublicData, error)
	}
	accountAdder interface {
		AddAccount(ctx context.Context, acc *deiz.ClinicianAccount) error
	}
	loginAllower interface {
		AllowLogin(ctx context.Context, loginCredentials deiz.Credentials) error
	}
)

func handlePostRegistration(allower loginAllower) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var f deiz.Credentials
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

func handlePostClinicianAccount(adder accountAdder) echo.HandlerFunc {
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

func handleGetClinicianAccount(getter accountDataGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		acc, err := getter.GetClinicianAccountData(ctx, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, acc)
	}
}

func handleGetClinicianAccountPublicData(getter accountDataGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID, err := strconv.Atoi(c.QueryParam("clinicianId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		data, err := getter.GetClinicianAccountPublicData(ctx, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, data)
	}
}
