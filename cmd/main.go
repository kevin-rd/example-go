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
	defer func(log *zap.Logger) {
		err := log.Sync()
		if err != nil {
			log.Info("Log shutdown failed", zap.Error(err))
		}
	}(log)

	addr := ":8080"
	router := server.InitRouter(log)
	srv := &http.Server{Addr: addr, Handler: router}

	go func() {
		log.Info("Http Server start ListenAndServe", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Http Server ListenAndServe failed", zap.Error(err))
		}
		log.Info("Http Server ListenAndServe stop.")
	}()

	stopSignal(log, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("HTTP Server shutdown failed", zap.Error(err))
		} else {
			log.Info("HTTP Server shutdown success")
		}
		return true
	})
}

// stopSignal watch stop signal and call call(), will exit program if call() return true
func stopSignal(log *zap.Logger, call func() bool) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	count := 0
	for sig := range stopCh {
		count++
		log.Warn("Receive signal", zap.String("signal", sig.String()), zap.Int("count", count))

		if count == 1 {
			log.Info("First signal received, initiating graceful shutdown...")
			go func() {
				if ok := call(); ok {
					log.Debug("call() ok, exiting.")
					os.Exit(0)
				}
			}()
		} else if count >= 2 {
			log.Warn("Receive signal again, force exit.")
			os.Exit(1)
		}
	}
}
