package main

import _ "go.uber.org/automaxprocs"

import (
	"context"
	"errors"
	"github.io/kevin-rd/demo-go/internal/logger"
	"github.io/kevin-rd/demo-go/internal/server"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// 初始化日志
	log := logger.InitLogger()
	defer func() {
		if err := log.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	addr := ":8080"
	router := server.InitRouter(log)
	srv := &http.Server{Addr: addr, Handler: router}

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
		log.Info("All goroutines finished..")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Info("Http Server start ListenAndServe", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Http Server ListenAndServe failed", zap.Error(err))
		}
		log.Info("Http Server ListenAndServe stop")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		waitStopSignal(log, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.Error("HTTP Server shutdown failed", zap.Error(err))
			} else {
				log.Info("HTTP Server shutdown success")
			}
		})
	}()
}

// waitStopSignal wait stop signal and call call(), will return if call() return
func waitStopSignal(log *zap.Logger, fns ...func()) {
	stopSigCh := make(chan os.Signal, 1)
	signal.Notify(stopSigCh, syscall.SIGINT, syscall.SIGTERM)

	canStopCh := make(chan struct{})
	count := 0
	for {
		select {
		case sig := <-stopSigCh:
			count++
			log.Warn("Receive signal", zap.String("signal", sig.String()), zap.Int("count", count))

			if count == 1 {
				log.Info("First signal received, initiating graceful shutdown...")
				go func() {
					for _, fn := range fns {
						fn()
					}
					close(canStopCh)
				}()
			} else if count >= 3 {
				log.Error("Receive signal again, force exit.")
				os.Exit(1)
			}
		case <-canStopCh:
			return
		}
	}
}
