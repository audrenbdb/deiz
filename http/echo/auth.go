package echo

import (
	"context"
	firebaseAuth "firebase.google.com/go/auth"
	"github.com/audrenbdb/deiz"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

type (
	reqMiddleware struct {
		credentialsGetter credentialsGetter
	}
	echoCtxCredentials struct {
		echo.Context
		credentials deiz.Credentials
	}

	credentialsGetter func(ctx context.Context, tokenID string) (deiz.Credentials, error)
)

//retrieve credentials in the current echo context
func getCredFromEchoCtx(c echo.Context) deiz.Credentials {
	credCtx, ok := c.(*echoCtxCredentials)
	if !ok {
		return deiz.Credentials{}
	}
	return credCtx.credentials
}

//getBearerToken attempts to read bearer token in the header request
func getBearerToken(header http.Header) (string, error) {
	tokenSplit := strings.Split(header.Get("Authorization"), "Bearer ")
	if len(tokenSplit) < 2 {
		return "", errParsingBearerToken
	}
	return tokenSplit[1], nil
}

func roleMW(getCredentials credentialsGetter, minRole deiz.Role) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenID, err := getBearerToken(c.Request().Header)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}
			cred, err := getCredentials(c.Request().Context(), tokenID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}
			if cred.Role < deiz.Role(minRole) {
				return c.JSON(http.StatusUnauthorized, errUnauthorizedRole.Error())
			}
			return next(&echoCtxCredentials{c, cred})
		}
	}
}

func publicMW(getCredentials credentialsGetter) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenID, err := getBearerToken(c.Request().Header)
			if err != nil {
				return next(&echoCtxCredentials{c, deiz.Credentials{}})
			}
			cred, err := getCredentials(c.Request().Context(), tokenID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}
			return next(&echoCtxCredentials{c, cred})
		}
	}
}

func FirebaseCredentialsGetter(client *firebaseAuth.Client) credentialsGetter {
	return func(ctx context.Context, tokenID string) (deiz.Credentials, error) {
		token, err := client.VerifyIDToken(ctx, tokenID)
		if err != nil {
			return deiz.Credentials{}, err
		}
		claims := token.Claims
		roleStr, ok := claims["role"]
		if !ok {
			return deiz.Credentials{}, errReadAuthClaims
		}
		roleFloat, ok := roleStr.(float64)
		if !ok || float64(int(roleFloat)) != roleFloat {
			return deiz.Credentials{}, errReadAuthClaims
		}

		userIDStr, ok := claims["userId"]
		if !ok {
			return deiz.Credentials{}, errReadAuthClaims
		}
		userIDFloat, ok := userIDStr.(float64)
		if !ok || float64(int(userIDFloat)) != userIDFloat {
			return deiz.Credentials{}, errReadAuthClaims
		}
		return deiz.Credentials{
			Role:   deiz.Role(int(roleFloat)),
			UserID: int(userIDFloat),
		}, nil
	}
}
