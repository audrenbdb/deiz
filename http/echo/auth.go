package echo

import (
	"context"
	firebaseAuth "firebase.google.com/go/auth"
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
		credentials credentials
	}
	credentials struct {
		userID int
		role   int
	}
	credentialsGetter func(ctx context.Context, tokenID string) (credentials, error)
)

//retrieve credentials in the current echo context
func getCredFromEchoCtx(c echo.Context) credentials {
	credCtx := c.(*echoCtxCredentials)
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

func roleMW(getCredentials credentialsGetter, minRole int) func(next echo.HandlerFunc) echo.HandlerFunc {
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
			if cred.role < minRole {
				return c.JSON(http.StatusUnauthorized, errUnauthorizedRole.Error())
			}
			return next(&echoCtxCredentials{c, cred})
		}
	}
}

func FirebaseCredentialsGetter(client *firebaseAuth.Client) credentialsGetter {
	return func(ctx context.Context, tokenID string) (credentials, error) {
		token, err := client.VerifyIDToken(ctx, tokenID)
		if err != nil {
			return credentials{}, err
		}
		claims := token.Claims
		roleStr, ok := claims["role"]
		if !ok {
			return credentials{}, errReadAuthClaims
		}
		roleFloat, ok := roleStr.(float64)
		if !ok || float64(int(roleFloat)) != roleFloat {
			return credentials{}, errReadAuthClaims
		}

		userIDStr, ok := claims["userId"]
		if !ok {
			return credentials{}, errReadAuthClaims
		}
		userIDFloat, ok := userIDStr.(float64)
		if !ok || float64(int(userIDFloat)) != userIDFloat {
			return credentials{}, errReadAuthClaims
		}
		return credentials{
			role:   int(roleFloat),
			userID: int(userIDFloat),
		}, nil
	}
}
