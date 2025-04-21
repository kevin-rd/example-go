package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const TimeFormatDataTimeMill = "2006-01-02 15:04:05.000-07"

// InitLogger 初始化 zap 日志
func InitLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(TimeFormatDataTimeMill))
	}
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeCaller = zapcore.FullCallerEncoder

	var encoder zapcore.Encoder
	var core zapcore.Core
	if isRunningInKubernetes() {
		encoder = zapcore.NewJSONEncoder(config)
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	}
	return zap.New(core)
}

func isRunningInKubernetes() bool {
	_, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	return exists
}
