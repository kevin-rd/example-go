package main

import (
	"context"
	"github.io/kevin-rd/demo-go/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: &app.ServerHandler{},
	}

	var mainWg sync.WaitGroup
	mainWg.Add(1)

	go func() {
		defer mainWg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP Server listen and serve failed: %v", err)
			os.Exit(0)
		}
		log.Println("Http Server has been Closed.")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	s := <-ch
	close(ch)
	log.Printf("Received signal: %v", s)
	if err := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(ctx)
	}(); err != nil {
		log.Printf("HTTP Server shutdown failed: %v", err)
	}
	mainWg.Wait()
	log.Println("Process has been Exit.")
}
