package booking2

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/auth"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type echoHandler struct {
	router         *echo.Echo
	getCredentials auth.GetCredentialsFromHTTPRequest
}

func NewEchoHandler(router *echo.Echo, getCredentials auth.GetCredentialsFromHTTPRequest) *echoHandler {
	return &echoHandler{router: router, getCredentials: getCredentials}
}

func (h *echoHandler) HandleGetBookings(
	clinicianGetBookings echo.HandlerFunc,
	patientGetBookings echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := h.getCredentials(c.Request())
		switch cred.Role {
		case deiz.ClinicianRole:
			return clinicianGetBookings(c)
		case deiz.PatientRole:
			return patientGetBookings(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}

func (h *echoHandler) HandlePatientGetBookings() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (h *echoHandler) HandleClinicianGetBookings(
	getClinicianCalendar echo.HandlerFunc,
	getClinicianUnpaidBookings echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.QueryParam("request")
		switch request {
		case "calendar":
			return getClinicianCalendar(c)
		case "unpaid":
			return getClinicianUnpaidBookings(c)
		}
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parameters")
	}
}

func (h *echoHandler) handleGetClinicianCalendar(
	getClinicianCalendar getClinicianCalendar,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred := h.getCredentials(c.Request())
		params, err := h.parseQueryCalendarParams(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		weekBookings, err := getClinicianCalendar(c.Request().Context(),
			params.from, params.bookingDuration, cred.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, weekBookings)
	}
}

type queryCalendarParams struct {
	from            time.Time
	bookingDuration time.Duration
}

func (h *echoHandler) parseQueryCalendarParams(c echo.Context) (queryCalendarParams, error) {
	duration, err := h.getIntegerQueryParam(c, "bookingDuration")
	if err != nil {
		return queryCalendarParams{}, err
	}
	from, err := h.getTimeFromUnixURLParam(c, "from")
	if err != nil {
		return queryCalendarParams{}, err
	}
	return queryCalendarParams{
		bookingDuration: time.Minute * time.Duration(duration),
		from:            from,
	}, nil
}

func (h *echoHandler) getTimeFromUnixURLParam(c echo.Context, paramName string) (time.Time, error) {
	i, err := strconv.ParseInt(c.QueryParam(paramName), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0).UTC(), nil
}

func (h *echoHandler) getIntegerQueryParam(c echo.Context, paramName string) (int, error) {
	return strconv.Atoi(c.QueryParam(paramName))
}

func (h *echoHandler) handleGetUnpaidBookings(c echo.Context, cred deiz.Credentials) error {
	return nil
}

func (h *echoHandler) handlePostBookings(auth auth.GetCredentialsFromHTTPRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		//cred := auth(c.Request())
		return nil
	}
}

func (h *echoHandler) handlePatchBookings(auth auth.GetCredentialsFromHTTPRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		//cred := auth(c.Request())
		return nil
	}
}

func (h *echoHandler) handleDeleteBookings(auth auth.GetCredentialsFromHTTPRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		//cred := auth(c.Request())
		return nil
	}
}
