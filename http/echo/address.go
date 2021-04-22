package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type (
	officeAddressAdder interface {
		AddClinicianOfficeAddress(ctx context.Context, address *deiz.Address, clinicianID int) error
	}
	homeAddressSetter interface {
		SetHomeAddress(ctx context.Context, address *deiz.Address, clinicianID int) error
	}
	addressDeleter interface {
		DeleteAddress(ctx context.Context, addressID, clinicianID int) error
	}
)

func handleDeleteClinicianAddress(deleter addressDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		addressID, err := strconv.Atoi(c.Param("aid"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = deleter.DeleteAddress(ctx, addressID, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostClinicianAddress(
	officeAddressAdder officeAddressAdder, homeAddressSetter homeAddressSetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var a deiz.Address
		if err = c.Bind(&a); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}

		//addressType can either be professional address or personal address
		addressType := c.QueryParam("type")
		if addressType != "office" && addressType != "home" {
			return c.JSON(http.StatusBadRequest, "address type not specified in the url")
		}

		if addressType == "home" {
			err = homeAddressSetter.SetHomeAddress(ctx, &a, clinicianID)
		} else {
			err = officeAddressAdder.AddClinicianOfficeAddress(ctx, &a, clinicianID)
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, a)
	}
}
