package echo

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
)

type (
	calendarSettingsEditer interface {
		EditCalendarSettings(ctx context.Context, s *deiz.CalendarSettings, clinicianID int) error
	}
)

func handlePatchCalendarSettings(editer calendarSettingsEditer) echo.HandlerFunc {
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
