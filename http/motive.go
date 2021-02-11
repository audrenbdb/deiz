package http

/*
func handlePostBookingMotive(addMotive deiz.AddBookingMotive, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var b deiz.BookingMotive
		if err := c.Bind(&b); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := validate.StructExceptCtx(ctx, b, "ID"); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := addMotive(ctx, &b, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, b)
	}
}

func handleDeleteBookingMotive(deleteMotive deiz.RemoveBookingMotive, validate validater) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		clinicianID := getCredFromEchoCtx(c).userID
		var m deiz.BookingMotive
		if err := c.Bind(&m); err != nil {
			return c.JSON(http.StatusBadRequest, errBind)
		}
		if err := validate.StructCtx(ctx, m); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := deleteMotive(ctx, &m, clinicianID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
}
*/
