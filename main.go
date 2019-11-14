package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	eventsFileName = "events.csv"
	port           = 5000
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

type Event struct {
	Latitude  float64
	Longitude float64
	Type      string
}

func run() error {
	events, err := extractEventsFromCSV(eventsFileName)
	if err != nil {
		log.Fatalln(err)
	}
	svr := newServer()
	svr.eventsDB = events

	log.Println("Listening on port : ", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), svr)
	if err != nil {
		return errors.Wrap(err, "Couldn't launch server")
	}
	return nil
}

func extractEventsFromCSV(filename string) ([]Event, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Problem opening file %q", filename))
	}

	csv := csv.NewReader(strings.NewReader(string(file)))
	records, err := csv.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read CSV file")
	}

	var events []Event

	for i, record := range records {
		if i == 0 {
			continue
		}
		latitude, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, errors.Wrap(err, "Couldn't convert to float")
		}
		longitude, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "Couldn't convert to float")
		}
		event := Event{
			Latitude:  latitude,
			Longitude: longitude,
			Type:      record[2],
		}
		events = append(events, event)
	}
	return events, nil
}
