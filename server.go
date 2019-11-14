package main

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	eventsDB []Event
	http.Handler
}

func newServer() *server {
	svr := &server{}

	router := mux.NewRouter()
	router.HandleFunc("/impressionsAndClicks", svr.handleImpressionsAndClicks).Methods(http.MethodPost)

	svr.Handler = router
	return svr
}

type POI struct {
	Name        string  `json:"name"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Impressions int     `json:"impressions"`
	Clicks      int     `json:"clicks"`
}

func (s *server) handleImpressionsAndClicks(w http.ResponseWriter, r *http.Request) {
	var POIs []POI

	err := json.NewDecoder(r.Body).Decode(&POIs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(POIs) != 2 {
		w.WriteHeader(http.StatusBadRequest)
	}

	POI0, POI1 := s.CalculateImpressionsAndClicks(POIs)

	result := map[string]POI{
		POI0.Name: POI0,
		POI1.Name: POI1,
	}

	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *server) CalculateImpressionsAndClicks(POIs []POI) (POI, POI) {
	for _, event := range s.eventsDB {
		if POIs[0].calculateSquaredDistanceFrom(event) > POIs[1].calculateSquaredDistanceFrom(event) {
			POIs[1].incrementImpressionOrClick(event.Type)
			continue
		}
		POIs[0].incrementImpressionOrClick(event.Type)
	}
	return POIs[0], POIs[1]
}

func (p *POI) incrementImpressionOrClick(eventType string) {
	if eventType == "click" {
		p.Clicks++
	} else if eventType == "imp" {
		p.Impressions++
	}
}

func (p POI) calculateSquaredDistanceFrom(e Event) float64 {
	return math.Pow(p.Lon-e.Longitude, 2) + math.Pow(p.Lat-e.Latitude, 2)
}
