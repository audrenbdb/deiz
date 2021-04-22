package account

import (
	"context"
	"github.com/audrenbdb/deiz"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockAccountCreater struct {
	err error
}

func (m *mockAccountCreater) CreateClinicianAccount(ctx context.Context, account *deiz.ClinicianAccount) error {
	return m.err
}

func TestAddAccount(t *testing.T) {
	//should fail to create an account
	usecase := AddAccountUsecase{
		&mockAccountCreater{err: deiz.GenericError},
	}
	errorOutput := deiz.GenericError

	err := usecase.AddAccount(context.Background(), &deiz.ClinicianAccount{})
	assert.Equal(t, errorOutput, err)
}
