package server

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, log *zap.Logger)

// HandleHello 处理 /hello 请求
func HandleHello(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
	// handle request
	var pi float64
	for i := 0; i < 10000; i++ {
		pi += (4.0 / (float64)(2.0*i+1.0)) * (1.0 - (2.0*(float64)(i%2))/1.0)
	}

	time.Sleep(time.Second * 5)

	log.Info("handle request: hello",
		zap.Float64("pi", pi),
		zap.String("log_tag", "module01"),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("client_ip", r.RemoteAddr),
		zap.Int8("Int8", 28),
	)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello!\n"))
}

func HandleHealth(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
