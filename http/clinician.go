package http

import (
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

func handlePatchClinicianEmail(edit deiz.EditClinicianEmail, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		type data struct {
			Email string `json:"email" validate:"email"`
		}
		var d data
		if err := c.Bind(&d); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if err := validate.StructCtx(ctx, d); err != nil {
			return c.JSON(http.StatusBadRequest, errValidating.Error())
		}
		err := edit(ctx, d.Email, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePatchClinicianPhone(edit deiz.EditClinicianPhone, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		type contact struct {
			Phone string `json:"phone" validate:"required,min=10"`
		}
		var f contact
		if err := c.Bind(&f); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if err := validate.StructCtx(ctx, f); err != nil {
			return c.JSON(http.StatusBadRequest, errValidating.Error())
		}
		err := edit(ctx, f.Phone, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePostClinician(add deiz.AddClinicianAccount, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var cl deiz.Clinician
		if err := c.Bind(&cl); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if err := validate.StructExceptCtx(ctx, cl, "ID"); err != nil {
			return c.JSON(http.StatusBadRequest, errValidating.Error())
		}
		err := add(ctx, &cl)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, cl)
	}
}
