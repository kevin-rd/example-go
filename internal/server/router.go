package server

import (
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.io/kevin-rd/demo-go/internal/metrics"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// InitRouter 初始化 HTTP 服务器的路由
func InitRouter(log *zap.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/hello", handleWrap(HandleHello, log))
	mux.HandleFunc("/basic-auth", handleWrap(basicAuthHandleWrap(HandleHello, "user_xxx", "passwd_xxx"), log))
	mux.HandleFunc("/token-auth", handleWrap(tokenAuthHandleWrap(HandleHello, "token_xxx"), log))

	return mux
}

func handleWrap(next HandlerFunc, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// pre handle
		start := time.Now()
		traceId := r.Header.Get("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		log = log.With(zap.String("trace_id", traceId))

		// handle
		next(w, r, log)

		// post handle
		duration := time.Since(start)
		metrics.RequestsCost.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
		if duration > time.Second*2 {
			log.Warn("slow request",
				zap.String("log_tag", "http_server"),
				zap.Duration("duration", duration),
				zap.String("url", r.URL.String()),
			)
		}
	}
}
