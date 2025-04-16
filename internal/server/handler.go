package server

import (
	"github.com/google/uuid"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HandleHello 处理 /hello 请求
func HandleHello(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
	start := time.Now()
	traceId := uuid.New().String()

	// 记录请求日志
	log.Info("receive request: hello",
		zap.String("trace_id", traceId),
		zap.String("log_tag", "module01"),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("client_ip", r.RemoteAddr),
		zap.Int8("Int8", 28),
	)

	// 处理请求
	value := r.URL.Query().Get("key")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, Zap!\n"))

	// 记录处理时间
	duration := time.Since(start)
	log.Info("handle complete",
		zap.String("trace_id", traceId),
		zap.String("log_tag", "module1"),
		zap.Duration("duration", duration),
		zap.String("key_test", value),
	)
}
