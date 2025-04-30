package main

import (
	"context"
	"errors"
	"github.io/kevin-rd/demo-go/internal/logger"
	"github.io/kevin-rd/demo-go/internal/server"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化日志
	log := logger.InitLogger()
	defer log.Sync()

	addr := ":8080"
	router := server.InitRouter(log)
	srv := &http.Server{Addr: addr, Handler: router}

	go func() {
		log.Info("Server start", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server failed", zap.Error(err))
		}
		log.Info("Server stop")
	}()

	stopSignal(log, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("HTTP Server shutdown failed", zap.Error(err))
		} else {
			log.Info("Process has been Exit.")
		}
	})
}

func stopSignal(log *zap.Logger, onClose func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	count := 0
	for sig := range stopCh {
		count++
		log.Warn("Receive signal", zap.String("signal", sig.String()), zap.Int("count", count))

		if count == 1 {
			log.Info("First signal received, initiating graceful shutdown...")
			go onClose()
		} else if count >= 3 {
			log.Warn("Receive signal again, force exit.")
			os.Exit(1)
		}
	}
}
