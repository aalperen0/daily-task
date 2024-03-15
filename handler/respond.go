package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, status int, msg string) {
	if status >= http.StatusInternalServerError {
		log.Fatalf("Responding with %d error %s", status, msg)
	}
	if msg == "" {
		msg = http.StatusText(status)
	}

	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}
	RespondWithJSON(w, http.StatusBadRequest, errorResponse)
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal json response %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)

}
