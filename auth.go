package deiz

import (
	"context"
)

//repo functions
type (
	logger interface {
		Login(ctx context.Context, email, password string) error
	}
)

//core functions
type (
	//Login logs an user in the application
	Login func(ctx context.Context, email, password string) error
)

func login(logger logger) Login {
	return func(ctx context.Context, email, password string) error {
		return logger.Login(ctx, email, password)
	}
}
