package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

type (
	clinicianPhoneEditer interface {
		EditClinicianPhone(ctx context.Context, phone string, clinicianID int) error
	}
	clinicianEmailEditer interface {
		EditClinicianEmail(ctx context.Context, email string, clinicianID int) error
	}
	addressEditer interface {
		EditAddress(ctx context.Context, address *deiz.Address, clinicianID int) error
	}
	clinicianAdeliEditer interface {
		EditClinicianAdeli(ctx context.Context, identifier string, clinicianID int) error
	}
	clinicianProfessionEditer interface {
		EditClinicianProfession(ctx context.Context, profession string, clinicianID int) error
	}
)

func handlePatchClinicianProfession(edit clinicianProfessionEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
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

func handlePatchClinicianPhone(edit clinicianPhoneEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
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

func handlePatchClinicianEmail(edit clinicianEmailEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
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

func handlePatchClinicianAddress(edit addressEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var a deiz.Address
		if err := c.Bind(&a); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := edit.EditAddress(ctx, &a, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}

func handlePatchClinicianAdeli(edit clinicianAdeliEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
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
