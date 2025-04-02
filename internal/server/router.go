package server

import (
	"net/http"

	"go.uber.org/zap"
)

// InitRouter 初始化 HTTP 服务器的路由
func InitRouter(log *zap.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		HandleHello(w, r, log)
	})
	return mux
}
