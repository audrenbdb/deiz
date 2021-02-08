package deiz_test

import (
	"context"
	"github.com/audrenbdb/deiz"
	"testing"
)

type mockAccountGetter struct {
	err error
}

func (c *mockAccountGetter) GetClinicianAccount(ctx context.Context, clinicianID int) (deiz.ClinicianAccount, error) {
	return deiz.ClinicianAccount{}, c.err
}

func TestGetClinicianAccount(t *testing.T) {
	/*
		r := deiz.Repo{}

		//should fail to get clinician account
		r.ClinicianAccount.Getter = &mockAccountGetter{err: errors.New("failed to get")}
		_, err := r.GetClinicianAccount(context.Background(), 1)
		assert.Error(t, err)

		//should succeed in getting clinicianAccount
		r.ClinicianAccount.Getter = &mockAccountGetter{}
		_, err = r.GetClinicianAccount(context.Background(), 1)
		assert.NoError(t, err)

	*/
}
