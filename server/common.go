package server

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
)

var errRequestorInvalid = errors.New("Requestor is anonymous or invalid")

func Requestor(key any, ctx context.Context) (r string, err error) {
	err = connect.NewError(connect.CodeUnauthenticated, errRequestorInvalid)

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
