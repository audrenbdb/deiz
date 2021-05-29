package booking2

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/auth"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type motive struct {
	name string
}

func handleGetBookings(auth auth.CredentialsFromHttpRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := auth(c.Request())
		switch cred.Role {
		case deiz.ClinicianRole:
			return handleClinicianGetBookings(c, cred)
		}
		return nil
	}
}

func handleClinicianGetBookings(c echo.Context, cred deiz.Credentials) error {
	request := c.QueryParam("request")
	switch request {
	case "calendar":
		return handleGetClinicianCalendar(c, cred)
	}
	//week calendar
	//or unpaid bookings

	return nil
}

func handleGetClinicianCalendar(c echo.Context, cred deiz.Credentials) error {
	from, err := getTimeFromParam(c, "from")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	motive, err := getMotiveFromParam(c)
	return nil
}

func getTimeFromParam(c echo.Context, paramName string) (time.Time, error) {
	i, err := strconv.ParseInt(c.QueryParam(paramName), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0).UTC(), nil
}

func handleGetUnpaidBookings(c echo.Context, cred deiz.Credentials) error {
	return nil
}

func handlePostBookings(auth auth.CredentialsFromHttpRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := auth(c.Request())
		return nil
	}
}

func handlePatchBookings(auth auth.CredentialsFromHttpRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := auth(c.Request())
		return nil
	}
}

func handleDeleteBookings(auth auth.CredentialsFromHttpRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := auth(c.Request())
		return nil
	}
}

func RegisterService(e *echo.Echo) {}
