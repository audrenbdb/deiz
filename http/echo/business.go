package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

type (
	businessEditer interface {
		EditClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error
	}
)

func handlePatchBusiness(patcher businessEditer) echo.HandlerFunc {
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
