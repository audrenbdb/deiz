package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

func handleDeleteClinicianAddress(deleter usecase.AddressDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		addressID, err := getURLIntegerParam(c, "aid")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = deleter.DeleteAddress(ctx, addressID, getCredFromEchoCtx(c))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostClinicianAddress(
	officeAddressAdder usecase.OfficeAddressAdder) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		ctx := c.Request().Context()
		credentials := getCredFromEchoCtx(c)
		a, err := getAddressFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		err = officeAddressAdder.AddClinicianOfficeAddress(ctx, &a, credentials)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, a)
	}
}

func getAddressFromRequest(c echo.Context) (deiz.Address, error) {
	var a deiz.Address
	return a, c.Bind(&a)
}
