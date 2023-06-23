package core

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var JWT_ACCOUNT_CLAIM string

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("JWT_ACCOUNT_CLAIM", "account")

	JWT_ACCOUNT_CLAIM = viper.GetString("JWT_ACCOUNT_CLAIM")
}

type ContextKey string

const ChatAccount = ContextKey("account")

func In[T comparable](o T, a []T) bool {
	for _, el := range a {
		if el == o {
			return true
		}
	}
	return false
}

func Unique[T comparable](a []T) []T {
	var b []T
	for _, el := range a {
		if !In(el, b) {
			b = append(b, el)
		}
	}
	return b
}

func NewLogger() (log *zap.Logger) {
	viper.SetDefault("LOG_LEVEL", int(zap.DebugLevel))
	level := viper.GetInt("LOG_LEVEL")

	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.Level(level))

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
}
