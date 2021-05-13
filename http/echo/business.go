package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
)

func handlePatchBusiness(patcher usecase.BusinessEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

		var b deiz.Business
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := patcher.EditClinicianBusiness(ctx, &b, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePatchBusinessAddress(patcher usecase.BusinessAddressEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

		a, err := getAddressFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = patcher.UpdateClinicianBusinessAddress(ctx, &a, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostBusinessAddress(poster usecase.BusinessAddressSetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID

		a, err := getAddressFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = poster.SetClinicianBusinessAddress(ctx, &a, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, &a)
	}
}
