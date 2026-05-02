package main

import (
	"net/http"
	"fmt"
	"log"
	"io"
	"encoding/json"
	"github.com/spf13/viper"
)

func initConfig() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    viper.AddConfigPath(".")
		viper.AddConfigPath("..")          

    // Read the config file
    if err := viper.ReadInConfig(); err != nil {
        // Handle the case where config file is not found
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Println("No config file found, using defaults")
            return nil
        }
        return fmt.Errorf("error reading config file: %w", err)
    }

    log.Printf("Using config file: %s", viper.ConfigFileUsed())
    return nil
}

func main() {
	stopNumbers := []int{
		1001694, // my house south 
		1001212, // farragut to north
  }
	url := "https://api.wmata.com/NextBusService.svc/json/jPredictions?StopID=%d"

  if err := initConfig(); err != nil {
  	log.Fatalf("Failed to initialize config: %v", err)
  }

	apiKey := viper.GetString("api.key")

	for _, stopNumber := range stopNumbers {
		request, _ := http.NewRequest("GET", fmt.Sprintf(url, stopNumber), nil)
		request.Header.Set("Cache-Control", "no-cache")
		request.Header.Set("api_key", apiKey)
	
		client := &http.Client{}
		response, _ := client.Do(request)

		defer response.Body.Close()

		// Read and print the response body
		body, _ := io.ReadAll(response.Body)

		var nbs NextBusServiceResponse
		err := json.Unmarshal(body, &nbs)
		if err != nil {
			fmt.Println("%v",err)
			return
		}

		for _, prediction := range nbs.Predictions {
			if prediction.RouteID == "D72" {
				fmt.Println(prediction)
			}
		}
	}


	// metro stations
	metroStations := []string{
		"N04",
		"C03",
	}

	for _, metroStation := range metroStations {
		url = "https://api.wmata.com/StationPrediction.svc/json/GetPrediction/%s"
 
		request, _ := http.NewRequest("GET", fmt.Sprintf(url, metroStation), nil)
		request.Header.Set("Cache-Control", "no-cache")
		request.Header.Set("api_key", apiKey)
	
		client := &http.Client{}
		response, _ := client.Do(request)

		defer response.Body.Close()
	
		body, _ := io.ReadAll(response.Body)

		var sp StationResponse 
		err := json.Unmarshal(body, &sp)
		if err != nil {
			fmt.Println("%v", err)
			return
		}

		for _, train := range sp.Trains {
			if train.Line == "SV" {
				fmt.Println(train)
			}
		} 
	}
}

type StationResponse struct {
	Trains []TrainPrediction `json:Trains`
}

type TrainPrediction struct {
	Cars string `json:"Car"`
	Destination string `json:"Destination"`
	DestinationCode string `json:"DestinationCode"`
	DestinationName string `json:"DestinationName"`
	Group string `json:"Group"`
	Line string `json:"Line"`
	LocationCode string `json:"LocationCode"`
	LocationName string `json:"LocationName"`
	Minutes string `json:"Min"`
}

type NextBusServiceResponse struct {
	StopName string `json:"StopName"`
	Predictions []Prediction `json:"Predictions"`
}

type Prediction struct {
	RouteID string `json:"RouteID"`
	DirectionText string `json:"DirectionText"`
	Minutes int `json:"Minutes"`
}


