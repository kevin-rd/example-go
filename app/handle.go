package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestData struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("method is not post")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read body error: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var requestData RequestData
	if err := json.Unmarshal(body, &requestData); err != nil {
		log.Printf("unmarshal body error: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Received request: %s", string(body))

	// handle
	response := Response{Message: fmt.Sprintf("Hello, %s!", requestData.Name)}
	json.NewEncoder(w).Encode(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
