package auth

import (
	"context"
	firebaseAuth "firebase.google.com/go/auth"
	"fmt"
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
	claims := token.Claims
	role, ok := claims["role"].(float64)
	if !ok {
		return deiz.Credentials{}
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		return deiz.Credentials{}
	}
	return deiz.Credentials{
		Role:   deiz.Role(int(role)),
		UserID: int(userID),
	}
}

func FirebaseHTTP(client *firebaseAuth.Client) GetCredentialsFromHttpRequest {
	return func(r *http.Request) deiz.Credentials {
		tokenID := getBearerTokenFromHTTPHeader(r.Header)
		if tokenID != "" {
			return getCredentialsFromFirebaseToken(r.Context(), client, tokenID)
		}
		return deiz.Credentials{UserID: 0, Role: 0}
	}
}

func MockHTTP(cred deiz.Credentials) GetCredentialsFromHttpRequest {
	return func(r *http.Request) deiz.Credentials {
		return cred
	}
}

type GetCredentialsFromHttpRequest = func(r *http.Request) deiz.Credentials

type fbAuth struct {
	client *firebaseAuth.Client
}

type CredentialsGetter interface {
	GetCredentialsFromJWT(ctx context.Context, jwt string) deiz.Credentials
}

type mockAuth struct {
	credentials deiz.Credentials
}

func (auth *mockAuth) GetCredentialsFromJWT(ctx context.Context, jwt string) deiz.Credentials {
	return auth.credentials
}

func (auth *fbAuth) GetCredentialsFromJWT(ctx context.Context, jwt string) deiz.Credentials {
	token, err := auth.client.VerifyIDToken(ctx, jwt)
	if err != nil {
		return deiz.Credentials{}
	}
	c, err := auth.getCredentialsFromToken(token)
	if err != nil {
		return deiz.Credentials{}
	}
	return c
}

func (auth *fbAuth) getCredentialsFromToken(token *firebaseAuth.Token) (deiz.Credentials, error) {
	return parseClaimsCredentials(token.Claims)
}

func parseClaimsCredentials(claims map[string]interface{}) (deiz.Credentials, error) {
	role, ok := claims["role"].(float64)
	if !ok {
		return deiz.Credentials{}, fmt.Errorf("not able to parse role claim")
	}
	userID, ok := claims["userId"].(float64)
	if !ok {
		return deiz.Credentials{}, fmt.Errorf("not able to parse userId claim")
	}
	return deiz.Credentials{
		Role:   deiz.Role(int(role)),
		UserID: int(userID),
	}, nil
}
