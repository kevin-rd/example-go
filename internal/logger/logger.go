package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// InitLogger 初始化 zap 日志
func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("fail to init zap logger: %v", err))
	}
	return logger
}
