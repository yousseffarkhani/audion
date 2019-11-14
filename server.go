package main

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"

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
		return
	}

	result := s.CalculateImpressionsAndClicks(POIs)

	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *server) CalculateImpressionsAndClicks(POIs []POI) map[string]POI {
	type indexAndDistance struct {
		index             int
		distanceFromEvent float64
	}
	for _, event := range s.eventsDB {
		var indexesAndDistances []indexAndDistance

		for index, POI := range POIs {
			indexesAndDistances = append(indexesAndDistances,
				indexAndDistance{
					index,
					POI.calculateSquaredDistanceFrom(event),
				},
			)
		}
		sort.Slice(indexesAndDistances, func(i, j int) bool {
			return indexesAndDistances[i].distanceFromEvent < indexesAndDistances[j].distanceFromEvent
		})

		POIs[indexesAndDistances[0].index].incrementImpressionOrClick(event.Type)
	}

	result := make(map[string]POI)

	for _, POI := range POIs {
		result[POI.Name] = POI
	}

	return result
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
