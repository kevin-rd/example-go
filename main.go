package main

import (
	"github.io/kevin-rd/demo-go/app"
	"log"
	"net/http"
	"os"
)

func main() {
	server := &http.Server{
		Addr:    ":12345",
		Handler: &app.ServerHandler{},
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("HTTP Server listen and serve failed: %v", err)
		os.Exit(1)
	}
	log.Println("Http Server has been Closed.")
}
