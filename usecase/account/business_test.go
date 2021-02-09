package account_test

import (
	"context"
	"errors"
	"github.com/audrenbdb/deiz"
	"github.com/audrenbdb/deiz/usecase/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	mockClinicianBusinessUpdater struct {
		err error
	}
)

func (s *mockClinicianBusinessUpdater) UpdateClinicianBusiness(ctx context.Context, b *deiz.Business, clinicianID int) error {
	return s.err
}

func TestEditClinicianBusiness(t *testing.T) {
	var tests = []struct {
		description string

		updater *mockClinicianBusinessUpdater

		inBusiness    *deiz.Business
		inClinicianID int

		outError error
	}{
		{
			description: "should fail to update business",
			updater:     &mockClinicianBusinessUpdater{err: errors.New("failed to update")},
			outError:    errors.New("failed to update"),
		},
		{
			description: "should succeed",
			updater:     &mockClinicianBusinessUpdater{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			u := account.Usecase{
				BusinessUpdater: test.updater,
			}
			err := u.EditClinicianBusiness(context.Background(), test.inBusiness, test.inClinicianID)
			assert.Equal(t, test.outError, err)
		})
	}
}
