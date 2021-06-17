package booking2_test

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/auth"
	"github.com/audrenbdb/deiz/booking2"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetBookings(t *testing.T) {
	tests := []struct {
		description          string
		role                 deiz.Role
		clinicianGetBookings echo.HandlerFunc
		patientGetBookings   echo.HandlerFunc

		expectedResponseCode int
	}{
		{
			description: "clinician should fail getting bookings with error 400",
			role:        deiz.ClinicianRole,

			clinicianGetBookings: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusBadRequest)
			},
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			description: "patient should fail getting bookings with error 500",
			role:        deiz.PatientRole,

			patientGetBookings: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusInternalServerError)
			},
			expectedResponseCode: http.StatusInternalServerError,
		},
		{
			description: "clinician should succeed getting bookings",
			role:        deiz.ClinicianRole,

			clinicianGetBookings: func(c echo.Context) error {
				return c.JSON(http.StatusOK, nil)
			},
			expectedResponseCode: http.StatusOK,
		},
		{
			description: "patient should succeed getting bookings",
			role:        deiz.PatientRole,

			patientGetBookings: func(c echo.Context) error {
				return c.JSON(http.StatusOK, nil)
			},
			expectedResponseCode: http.StatusOK,
		},
		{
			description:          "other role should get unauthorized",
			role:                 deiz.PublicRole,
			expectedResponseCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		h := booking2.NewEchoHandler(echo.New(),
			auth.MakeMockCredentialsGetter(deiz.Credentials{UserID: 1, Role: test.role}),
		)
		e := echo.New()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		c := e.NewContext(req, rec)
		getBookings := h.HandleGetBookings(test.clinicianGetBookings, test.patientGetBookings)
		err := getBookings(c)
		if err != nil {
			he := err.(*echo.HTTPError)
			assert.Equal(t, test.expectedResponseCode, he.Code)
		} else {
			assert.Equal(t, test.expectedResponseCode, rec.Code)
		}
	}
}

func TestHandleClinicianGetBookings(t *testing.T) {
	tests := []struct {
		description string

		getClinicianCalendar       echo.HandlerFunc
		getClinicianUnpaidBookings echo.HandlerFunc
		req                        *http.Request
		expectedResponseCode       int
	}{
		{
			description:          "should fail to decode request param and return 400",
			req:                  newRequest("GET", "/", nil, nil),
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			description: "should fail to get clinician calendar with error 500",
			req:         newRequest("GET", "/", nil, map[string]string{"request": "calendar"}),
			getClinicianCalendar: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusInternalServerError)
			},
			expectedResponseCode: http.StatusInternalServerError,
		},
		{
			description: "should fail to get clinician unpaid bookings with error 500",
			req:         newRequest("GET", "/", nil, map[string]string{"request": "unpaid"}),
			getClinicianUnpaidBookings: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusInternalServerError)
			},
			expectedResponseCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		e := echo.New()
		h := booking2.NewEchoHandler(e,
			auth.MakeMockCredentialsGetter(deiz.Credentials{}),
		)
		rec := httptest.NewRecorder()
		c := e.NewContext(test.req, rec)
		clinicianGetBookings := h.HandleClinicianGetBookings(test.getClinicianCalendar, test.getClinicianUnpaidBookings)
		err := clinicianGetBookings(c)
		if err != nil {
			he := err.(*echo.HTTPError)
			assert.Equal(t, test.expectedResponseCode, he.Code)
		} else {
			assert.Equal(t, test.expectedResponseCode, rec.Code)
		}
	}
}

func mockHandleGetClinicianBookings(statusCode int, body interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(statusCode, body)
	}
}

func newRequest(method string, target string, body io.Reader, params map[string]string) *http.Request {
	req := httptest.NewRequest(method, target, body)
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req
}
