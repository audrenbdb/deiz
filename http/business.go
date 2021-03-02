package http

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

type (
	TaxExemptionCodesGetter interface {
		GetTaxExemptionCodes(ctx context.Context) ([]deiz.TaxExemption, error)
	}
	BusinessEditer interface {
		EditClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error
	}
)

func handleGetTaxExemptionCodes(getter TaxExemptionCodesGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		codes, err := getter.GetTaxExemptionCodes(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, codes)
	}
}

func handlePatchBusiness(patcher BusinessEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID

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
