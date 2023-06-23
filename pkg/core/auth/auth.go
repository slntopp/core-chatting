package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt"
	"github.com/slntopp/core-chatting/pkg/core"
	"go.uber.org/zap"
)

func NewUserInterceptor(log *zap.Logger, signing_key []byte) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			header := req.Header().Get("Authorization")
			log.Debug("Authorization Header", zap.String("header", header))

			segments := strings.Split(header, " ")
			if len(segments) != 2 {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
			}

			log.Debug("Validating Token", zap.String("token", segments[1]))
			claims, err := validateToken(signing_key, segments[1])

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

func validateToken(signing_key []byte, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return signing_key, nil
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
