package http

import (
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

func handlePostClinicianAddress(
	addPersonalAddress deiz.AddClinicianPersonalAddress, addOfficeAddress deiz.AddClinicianOfficeAddress,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var a deiz.Address
		if err = c.Bind(&a); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if !a.IsValid() {
			return c.JSON(http.StatusBadRequest, errValidating.Error())
		}

		//addressType can either be professional address or personal address
		addressType := c.QueryParam("addressType")
		if addressType != "professional" && addressType != "personal" {
			return c.JSON(http.StatusBadRequest, "address type not specified in the url")
		}

		if addressType == "personal" {
			err = addPersonalAddress(ctx, &a, clinicianID)
		} else {
			err = addOfficeAddress(ctx, &a, clinicianID)
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, a)
	}
}

func handlePatchClinicianAddress(updateAddress deiz.EditClinicianAddress, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var address deiz.Address
		if err := c.Bind(&address); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if err := validate.StructCtx(ctx, address); err != nil {
			return c.JSON(http.StatusBadRequest, errValidating)
		}
		err := updateAddress(ctx, &address, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
