package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrewkav/viber"
)

type statsResponse struct {
	Total int `json:"total"`
}

func handleStats(_ poll, _ *viber.Viber, s *storage, w http.ResponseWriter, _ *http.Request) {
	count, err := s.PersistCount()
	if err != nil {
		log.Printf("Error reading count: %v", err)
		http.Error(w, "Error reading count", http.StatusBadRequest)
		return
	}

	resp := statsResponse{
		Total: count,
	}

	json, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
		http.Error(w, "Error marshalling json", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
