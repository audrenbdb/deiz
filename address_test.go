package deiz_test

import (
	"context"
	"github.com/audrenbdb/deiz"
)

type (
	mockAddressUpdater struct {
		err error
	}
)

func (c *mockAddressUpdater) UpdateAddress(ctx context.Context, address *deiz.Address) error {
	return c.err
}
