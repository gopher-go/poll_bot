package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrewkav/viber"
)

type statsResponse struct {
	Total int `json:"total"`
	Aggs  struct {
		CountByResidenceLocation     map[string]int `json:"countByResidenceLocation"`
		CountByResidenceLocationType map[string]int `json:"countByResidenceLocationType"`
	} `json:"aggs"`
}

func handleStats(_ poll, _ *viber.Viber, s *storage, statsDao *statsDao, w http.ResponseWriter, r *http.Request) {
	var statsRespBytes []byte

	var err error
	count, err := s.CountCached()
	if err != nil {
		log.Printf("Error reading count: %v", err)
		http.Error(w, "Error reading count", http.StatusBadRequest)
		return
	}

	statResp := statsResponse{
		Total: count,
	}

	if statsDao != nil {
		aggs, err := statsDao.CountByFieldCached(r.Context(), residenceLocationType, residenceLocation)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Printf("unable to get aggregation by residence location type, err=%v", err)
			http.Error(w, "Unable to get statistic", http.StatusInternalServerError)
			return
		}

		statResp.Aggs.CountByResidenceLocationType = aggs[string(residenceLocationType)]
		statResp.Aggs.CountByResidenceLocation = aggs[string(residenceLocation)]
	}

	statsRespBytes, err = json.Marshal(statResp)
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
		http.Error(w, "Error marshalling json", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(statsRespBytes)
}
