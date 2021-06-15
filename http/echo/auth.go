package echo

import (
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	echoCtxCredentials struct {
		echo.Context
		credentials deiz.Credentials
	}
)

//retrieve credentials in the current echo context
func getCredFromEchoCtx(c echo.Context) deiz.Credentials {
	credCtx, ok := c.(*echoCtxCredentials)
	if !ok {
		return deiz.Credentials{}
	}
	return credCtx.credentials
}

func roleMW(auth auth.GetCredentialsFromHttpRequest, minRole deiz.Role) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cred := auth(c.Request())
			if cred.Role < deiz.Role(minRole) {
				return c.JSON(http.StatusUnauthorized, errUnauthorizedRole.Error())
			}
			return next(&echoCtxCredentials{c, cred})
		}
	}
}
