package server

import (
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.io/kevin-rd/demo-go/internal/metrics"
	"net/http"
	"strconv"
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

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
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
		log.Info("[http] receive request",
			zap.String("trace_id", traceId),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		// wrap ResponseWriter
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		// handle
		next(rec, r, log)

		// post handle
		cost := time.Since(start)
		metrics.RequestsCost.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rec.status)).Observe(cost.Seconds())

		// record warn log if cost > 2s
		if cost > time.Second*2 {
			log.Warn("[http] slow request",
				zap.String("log_tag", "http_server"),
				zap.Duration("cost", cost),
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
