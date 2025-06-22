package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Location struct {
	Count    int32  `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getLocations(url string) Location {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Fatal(err)
	}
	return location
}
