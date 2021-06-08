package auth

import (
	"context"
	firebaseAuth "firebase.google.com/go/auth"
	"github.com/audrenbdb/deiz"
	"net/http"
	"strings"
)

func getBearerTokenFromHTTPHeader(header http.Header) string {
	tokenSplit := strings.Split(header.Get("Authorization"), "Bearer ")
	if len(tokenSplit) < 2 {
		return ""
	}
	return tokenSplit[1]
}

func getCredentialsFromFirebaseToken(ctx context.Context, client *firebaseAuth.Client, tokenID string) deiz.Credentials {
	token, err := client.VerifyIDToken(ctx, tokenID)
	if err != nil {
		return deiz.Credentials{UserID: 0, Role: 0}
	}
	return getCredentialsFromTokenClaims(token.Claims)
}

func getCredentialsFromTokenClaims(claims map[string]interface{}) deiz.Credentials {
	role := parseIntClaim(claims["role"])
	userID := parseIntClaim(claims["userId"])
	return deiz.Credentials{
		Role:   deiz.Role(role),
		UserID: userID,
	}
}

func parseIntClaim(val interface{}) int {
	floatVal, ok := val.(float64)
	if !ok {
		return 0
	}
	return int(floatVal)

}

func CreateHTTPCredentialsGetterWithFirebase(client *firebaseAuth.Client) GetCredentialsFromHttpRequest {
	return func(r *http.Request) deiz.Credentials {
		tokenID := getBearerTokenFromHTTPHeader(r.Header)
		if tokenID != "" {
			return getCredentialsFromFirebaseToken(r.Context(), client, tokenID)
		}
		return deiz.Credentials{UserID: 0, Role: 0}
	}
}

func CreateHTTPCredentialsGetterWithMock(cred deiz.Credentials) GetCredentialsFromHttpRequest {
	return func(r *http.Request) deiz.Credentials {
		return cred
	}
}

type GetCredentialsFromHttpRequest = func(r *http.Request) deiz.Credentials
