package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt"
	"github.com/slntopp/core-chatting/pkg/core"
	"go.uber.org/zap"
)

type interceptor struct {
	log         *zap.Logger
	signing_key []byte
}

func NewAuthInterceptor(log *zap.Logger, signing_key []byte) *interceptor {
	return &interceptor{
		log:         log.Named("AuthInterceptor"),
		signing_key: signing_key,
	}
}

func (i *interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {

	return connect.UnaryFunc(func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		header := req.Header().Get("Authorization")
		i.log.Debug("Authorization Header", zap.String("header", header))

		segments := strings.Split(header, " ")
		if len(segments) != 2 {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
		}

		i.log.Debug("Validating Token", zap.String("token", segments[1]))
		claims, err := validateToken(i.signing_key, segments[1])

		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		acc := claims[core.JWT_ACCOUNT_CLAIM]
		ctx = context.WithValue(ctx, core.ChatAccount, acc.(string))

		return next(ctx, req)
	})
}

func (i *interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	i.log.Debug("WrapStreamingClient")
	return next
}
func (i *interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	i.log.Debug("Setup Wrap Streaming Handler")
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		header := shc.RequestHeader().Get("Authorization")
		i.log.Debug("Authorization Header", zap.String("header", header))

		segments := strings.Split(header, " ")
		if len(segments) != 2 {
			return connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
		}

		i.log.Debug("Validating Token", zap.String("token", segments[1]))
		claims, err := validateToken(i.signing_key, segments[1])

		if err != nil {
			return connect.NewError(connect.CodeUnauthenticated, err)
		}

		acc := claims[core.JWT_ACCOUNT_CLAIM]
		ctx = context.WithValue(ctx, core.ChatAccount, acc.(string))

		return next(ctx, shc)
	}
}

func (i *interceptor) WrapSocket(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		segments := strings.Split(header, " ")
		if len(segments) != 2 {
			return
		}

		i.log.Debug("Validating Token", zap.String("token", segments[1]))
		claims, err := validateToken(i.signing_key, segments[1])

		if err != nil {
			return
		}

		acc := claims[core.JWT_ACCOUNT_CLAIM]
		r.Header.Set(string(core.ChatAccount), acc.(string))

		next.ServeHTTP(w, r)
	})
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
