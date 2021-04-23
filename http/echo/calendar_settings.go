package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase"
	"github.com/labstack/echo"
	"net/http"
)

func handlePatchCalendarSettings(editer usecase.CalendarSettingsEditer) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var s deiz.CalendarSettings
		if err := c.Bind(&s); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		err := editer.EditCalendarSettings(ctx, &s, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
