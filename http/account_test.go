package http

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandleGetClinicianAccount(t *testing.T) {
	coreFail := func(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
		return deiz.ClinicianAccount{}, errors.New("test")
	}
	coreSuccess := func(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
		return deiz.ClinicianAccount{}, nil
	}
	goodRequest := fakeEchoRequest{queryParam: "clinicianId", queryParamValue: "7"}
	h := handleGetClinicianAccount(coreFail)
	//should fail to decode id
	c, rec := createFakeEchoRequest(fakeEchoRequest{})
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
	//should return internal server error
	c, rec = createFakeEchoRequest(fakeEchoRequest{queryParam: "clinicianId", queryParamValue: "7"})
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	h = handleGetClinicianAccount(coreSuccess)
	//should succeed
	c, rec = createFakeEchoRequest(goodRequest)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
