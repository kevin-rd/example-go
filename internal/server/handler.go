package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HandleHello 处理 /hello 请求
func HandleHello(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
	start := time.Now()

	// 记录请求日志
	log.Info("receive request",
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("client_ip", r.RemoteAddr),
		zap.Int8("Int8", 28),
	)

	// 处理请求
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, Zap!\n"))

	// 记录处理时间
	duration := time.Since(start)
	log.Info("handle complete",
		zap.Duration("duration", duration),
	)
}
