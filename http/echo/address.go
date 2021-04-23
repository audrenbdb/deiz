package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func handleDeleteClinicianAddress(deleter usecase.AddressDeleter) echo.HandlerFunc {
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
	officeAddressAdder usecase.OfficeAddressAdder, homeAddressSetter usecase.HomeAddressSetter) echo.HandlerFunc {
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
