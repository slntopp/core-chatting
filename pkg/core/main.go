package core

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
