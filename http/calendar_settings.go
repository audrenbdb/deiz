package http

/*
func handlePatchCalendarSettings(update deiz.EditCalendarSettings, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var s deiz.CalendarSettings
		if err := c.Bind(&s); err != nil {
			return c.JSON(http.StatusBadRequest, errBind.Error())
		}
		if err := validate.StructCtx(ctx, s); err != nil {
			return c.JSON(http.StatusBadRequest, errValidating.Error())
		}
		err := update(ctx, &s, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
*/
