package http

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type (
	ClinicianAccountGetter interface {
		GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error)
	}
	ClinicianAccountPublicDataGetter interface {
		GetClinicianAccountPublicData(ctx context.Context, clinicianID int) (deiz.ClinicianAccountPublicData, error)
	}
	ClinicianAccountAdder interface {
		AddClinicianAccount(ctx context.Context, acc *deiz.ClinicianAccount) error
	}
	ClinicianRegistrationCompleter interface {
		EnsureClinicianRegistrationComplete(ctx context.Context, email, password string) error
	}
)

func handlePostRegistration(completer ClinicianRegistrationCompleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		type reg struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var f reg
		if err := c.Bind(&f); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := completer.EnsureClinicianRegistrationComplete(ctx, f.Email, f.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostClinicianAccount(adder ClinicianAccountAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var acc deiz.ClinicianAccount
		if err := c.Bind(&acc); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := adder.AddClinicianAccount(ctx, &acc); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handleGetClinicianAccount(getter ClinicianAccountGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		acc, err := getter.GetClinicianAccount(ctx, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, acc)
	}
}

func handleGetClinicianAccountPublicData(getter ClinicianAccountPublicDataGetter) echo.HandlerFunc {
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
