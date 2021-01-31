package http

import "context"

type validater interface {
	StructCtx(ctx context.Context, s interface{}) (err error)
	StructExceptCtx(ctx context.Context, s interface{}, fields ...string) (err error)
}
