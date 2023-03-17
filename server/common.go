package server

import (
	"context"
	"errors"

	c "github.com/bufbuild/connect-go"
)

func Requestor(key any, ctx context.Context) (r string, err error) {
	err = NewError(c.CodeUnauthenticated, "Requestor is anonymous or invalid")

	v := ctx.Value(key)
	if v == nil {
		return
	}

	r, ok := v.(string)
	if !ok {
		return
	}
	return r, nil
}

func NewError(code c.Code, s string) error {
	return c.NewError(code, errors.New(s))
}
