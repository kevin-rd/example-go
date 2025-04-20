package server

import (
	"go.uber.org/zap"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, log *zap.Logger, traceId string)

// HandleHello 处理 /hello 请求
func HandleHello(w http.ResponseWriter, r *http.Request, log *zap.Logger, traceId string) {
	// handle request
	result := 3.14 * 8

	log.Info("receive request: hello",
		zap.Float64("result", result),
		zap.String("trace_id", traceId),
		zap.String("log_tag", "module01"),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("client_ip", r.RemoteAddr),
		zap.Int8("Int8", 28),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, Zap!\n"))
}
