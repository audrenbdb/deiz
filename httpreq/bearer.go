package http

import (
	"net/http"
	"strings"
)

func GetAuthBearerJWT(r *http.Request) string {
	tokenSplit := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(tokenSplit) < 2 {
		return ""
	}
	return tokenSplit[1]
}
