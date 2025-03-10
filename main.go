package main

import (
	"context"
	"errors"
	"github.io/kevin-rd/demo-go/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const port = 8080

func main() {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: &app.ServerHandler{},
	}

	var mainWg sync.WaitGroup
	mainWg.Add(1)

	go func() {
		defer mainWg.Done()
		log.Printf("Http Server listen and serve on %d", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
