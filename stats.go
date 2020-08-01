package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrewkav/viber"
	"github.com/coocood/freecache"
)

type statsResponse struct {
	Total int `json:"total"`
}

var (
	statsCacheKey = []byte("stats")
	statsCache    = freecache.NewCache(512 * 1024 * 1024)
)

func handleStats(_ poll, _ *viber.Viber, s *storage, w http.ResponseWriter, _ *http.Request) {
	var statsRespBytes []byte

	var err error
	statsRespBytes, err = statsCache.Get(statsCacheKey)
	if err != nil {
		count, err := s.PersistCount()
		if err != nil {
			log.Printf("Error reading count: %v", err)
			http.Error(w, "Error reading count", http.StatusBadRequest)
			return
		}

		statsRespBytes, err = json.Marshal(statsResponse{
			Total: count,
		})
		if err != nil {
			log.Printf("Error marshalling json: %v", err)
			http.Error(w, "Error marshalling json", http.StatusBadRequest)
			return
		}

		_ = statsCache.Set(statsCacheKey, statsRespBytes, 5)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(statsRespBytes)
}
