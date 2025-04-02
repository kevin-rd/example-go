package main

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.io/kevin-rd/demo-go/internal/logger"
	"github.io/kevin-rd/demo-go/internal/server"
)

func main() {
	// 初始化日志
	log := logger.InitLogger()
	defer log.Sync()

	router := server.InitRouter(log)
	srv := &http.Server{Addr: ":8080", Handler: router}
	go func() {
		addr := ":8080"
		log.Info("server start", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server start failed", zap.Error(err))
		}
		log.Info("server stop")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	s := <-ch
	close(ch)
	log.Info("Received signal", zap.String("signal", s.String()))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("HTTP Server shutdown failed", zap.Error(err))
	} else {
		log.Info("Process has been Exit.")
	}
}
