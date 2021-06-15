package http_test

import (
	"github.com/audrenbdb/deiz/http"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	r.Header.Set("Authorization", "Bearer toto")
	assert.Equal(t, http.GetAuthBearerJWT(r), "toto")

	r.Header.Set("Authorization", "")
	assert.Equal(t, http.GetAuthBearerJWT(r), "")
}
