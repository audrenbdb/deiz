package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
)

func handlePatchClinicianProfession(edit usecase.ClinicianProfessionEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		var post struct {
			Profession string `json:"profession"`
		}
		if err := c.Bind(&post); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := edit.EditClinicianProfession(ctx, post.Profession, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePatchClinicianPhone(edit usecase.ClinicianPhoneEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		var post struct {
			Phone string `json:"phone"`
		}
		if err := c.Bind(&post); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := edit.EditClinicianPhone(ctx, post.Phone, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePatchClinicianEmail(edit usecase.ClinicianEmailEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		type data struct {
			Email string `json:"email"`
		}
		var d data
		if err := c.Bind(&d); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		err := edit.EditClinicianEmail(ctx, d.Email, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePatchClinicianAddress(edit usecase.AddressEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var a deiz.Address
		if err := c.Bind(&a); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := edit.EditAddress(ctx, &a, getCredFromEchoCtx(c))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePatchClinicianAdeli(edit usecase.ClinicianAdeliEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).UserID
		var a deiz.Adeli
		if err := c.Bind(&a); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := edit.EditClinicianAdeli(ctx, a.Identifier, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
