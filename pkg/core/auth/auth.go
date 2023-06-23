package auth

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt"
	"github.com/slntopp/core-chatting/pkg/core"
)

var SIGNING_KEY []byte

func NewUserInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			token := req.Header().Get("bearer")

			claims, err := validateToken(token)

			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			acc := claims[core.JWT_ACCOUNT_CLAIM]
			ctx = context.WithValue(ctx, core.ChatAccount, acc.(string))

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return SIGNING_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, nil
}
