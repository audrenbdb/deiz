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
	if err == nil {
		claims := token.Claims
		role := claims["role"].(int)
		userID := claims["userId"].(int)
		return deiz.Credentials{UserID: userID, Role: deiz.Role(role)}
	}
	return deiz.Credentials{UserID: 0, Role: 0}
}

func FirebaseHTTP(client *firebaseAuth.Client) CredentialsFromHttpRequest {
	return func(r *http.Request) deiz.Credentials {
		tokenID := getBearerTokenFromHTTPHeader(r.Header)
		if tokenID != "" {
			return getCredentialsFromFirebaseToken(r.Context(), client, tokenID)
		}
		return deiz.Credentials{UserID: 0, Role: 0}
	}
}

func MockHTTP(cred deiz.Credentials) CredentialsFromHttpRequest {
	return func(r *http.Request) deiz.Credentials {
		return cred
	}
}

type CredentialsFromHttpRequest = func(r *http.Request) deiz.Credentials
