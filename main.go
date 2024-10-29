package main

import (
	"log"
	"net/http"
)

func doHello(w http.ResponseWriter, r *http.Request) {
	log.Println("do hello")
}

func main() {
	http.HandleFunc("/hello", doHello)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		return
	}
}
