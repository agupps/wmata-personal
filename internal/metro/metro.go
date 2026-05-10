package metro

import (
	"net/http"
	"fmt"
	"log"
	"io"
	"encoding/json"
	"github.com/spf13/viper"
	"wmata/internal/transit"
)

func MetroStops(metroStop string) string {
	// metro stations
  
	apiKey := viper.GetString("api.key")

	url := "https://api.wmata.com/StationPrediction.svc/json/GetPrediction/%s"

	request, _ := http.NewRequest("GET", fmt.Sprintf(url, metroStop), nil)
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("api_key", apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("failed to hit API: %v", err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var sp StationResponse 
	err = json.Unmarshal(body, &sp)
	if err != nil {
		fmt.Println("%v", err)
		return ""
	}

	silverLineTrainPredictions := []transit.Prediction{}

	for _, train := range sp.Trains {
		if train.Line == "SV" {
			silverLineTrainPredictions = append(silverLineTrainPredictions, train)
		}
	}
	return transit.PredictionsToString(silverLineTrainPredictions)
}

type StationResponse struct {
	Trains []*TrainPrediction `json:Trains`
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

func (t *TrainPrediction) GetMinutes() string {
	return t.Minutes
}

func (t *TrainPrediction) GetDirectionText() string {
	return t.DestinationName
}

func (t *TrainPrediction) GetRouteID() string {
	return t.Line
}

func (t TrainPrediction) GetPrediction() string {
  return fmt.Sprintf("%s to %s: %s\n", t.Line, t.DestinationName, t.Minutes)
}

