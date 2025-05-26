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

	// metrics
	mux.Handle("/metrics", promhttp.Handler())

	// health
	mux.HandleFunc("/health/startup", handleWrap(HandleHealth, log))
	mux.HandleFunc("/health/liveness", handleWrap(HandleHealth, log))
	mux.HandleFunc("/health/readiness", handleWrap(HandleHealth, log))

	// demo
	mux.HandleFunc("/hello", handleWrap(HandleHello, log))
	mux.HandleFunc("/basic-auth", basicAuthHandleWrap(HandleHello, log, "user_xxx", "passwd_xxx"))
	mux.HandleFunc("/token-auth", tokenAuthHandleWrap(HandleHello, log, "token_xxx"))
	mux.HandleFunc("/token-auth2", handleWrap(tokenAuthWrap(HandleHello, "token_xxx"), log))

	return mux
}

// handleWrap is a higher-order function
// that wraps a handler function HandlerFunc and adds logging capabilities before and after its execution.
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

// basicAuthHandleWrap is a higher-order function
// that wraps a handler function HandlerFunc and adds basic authentication capabilities before its execution.
func basicAuthHandleWrap(next HandlerFunc, log *zap.Logger, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleWrap(basicAuthWrap(next, username, password), log)
	}
}

// tokenAuthHandleWrap is a higher-order function
// that wraps a handler function HandlerFunc and adds token authentication capabilities before its execution.
func tokenAuthHandleWrap(next HandlerFunc, log *zap.Logger, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleWrap(tokenAuthWrap(next, token), log)
	}
}
