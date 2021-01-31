package http

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type fakeEchoRequest struct {
	body            string
	path            string
	param           string
	paramValue      string
	queryParam      string
	queryParamValue string
}

//mock a new request with given json body
func createFakeEchoRequest(r fakeEchoRequest) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	q := make(url.Values)
	if r.queryParam != "" {
		q.Set(r.queryParam, r.queryParamValue)
	}
	req := httptest.NewRequest(http.MethodPost, "/?"+q.Encode(), strings.NewReader(r.body))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if r.path != "" {
		c.SetPath(r.path)
	}
	if r.param != "" {
		c.SetParamNames(r.param)
		c.SetParamValues(r.paramValue)
	}
	return &echoCtxCredentials{c, credentials{userID: 7, role: 2}}, rec
}
