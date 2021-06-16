package booking2

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/auth"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type echoServer struct {
	router         *echo.Echo
	getCredentials auth.GetCredentialsFromHttpRequest
}

func NewEchoHandler(router *echo.Echo, getCredentials auth.GetCredentialsFromHttpRequest) *echoServer {
	return &echoServer{getCredentials: getCredentials}
}

func (e *echoServer) handleGetBookings(
	handleGetClinicianBookings echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := e.getCredentials(c.Request())
		switch cred.Role {
		case deiz.ClinicianRole:
			return handleGetClinicianBookings(c)
		}
		return nil
	}
}

func (e *echoServer) handleGetClinicianBookings(
	handleGetClinicianWeek echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.QueryParam("request")
		switch request {
		case "week":
			return handleGetClinicianWeek(c)
		case "unpaid":
			return nil
		}
		return nil
	}
}

func (e *echoServer) handleGetClinicianWeek(
	getClinicianWeek getClinicianWeek,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := e.getCredentials(c.Request())
		params, err := e.parseQueryWeekParams(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		weekBookings, err := getClinicianWeek(c.Request().Context(),
			params.from, params.bookingDuration, cred.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, weekBookings)
	}
}

type queryWeekParams struct {
	from            time.Time
	bookingDuration time.Duration
}

func (e *echoServer) parseQueryWeekParams(c echo.Context) (queryWeekParams, error) {
	duration, err := e.getIntegerQueryParam(c, "bookingDuration")
	if err != nil {
		return queryWeekParams{}, err
	}
	from, err := e.getTimeFromUnixURLParam(c, "from")
	if err != nil {
		return queryWeekParams{}, err
	}
	return queryWeekParams{
		bookingDuration: time.Minute * time.Duration(duration),
		from:            from,
	}, nil
}

func (e *echoServer) getTimeFromUnixURLParam(c echo.Context, paramName string) (time.Time, error) {
	i, err := strconv.ParseInt(c.QueryParam(paramName), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0).UTC(), nil
}

func (e *echoServer) getIntegerQueryParam(c echo.Context, paramName string) (int, error) {
	return strconv.Atoi(c.QueryParam(paramName))
}

func (e *echoServer) handleGetUnpaidBookings(c echo.Context, cred deiz.Credentials) error {
	return nil
}

func (e *echoServer) handlePostBookings(auth auth.GetCredentialsFromHttpRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		//cred := auth(c.Request())
		return nil
	}
}

func (e *echoServer) handlePatchBookings(auth auth.GetCredentialsFromHttpRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		//cred := auth(c.Request())
		return nil
	}
}

func (e *echoServer) handleDeleteBookings(auth auth.GetCredentialsFromHttpRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		//cred := auth(c.Request())
		return nil
	}
}
