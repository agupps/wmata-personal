package buses

import (
	"net/http"
	"fmt"
	"log"
	"io"
	"encoding/json"
	"github.com/samber/lo"
	"wmata/internal/transit"
	"wmata/internal/config"
)


func BusStops(stopNumber int, config *config.Config, busLineFilter *string) []transit.Prediction {
	url := "https://api.wmata.com/NextBusService.svc/json/jPredictions?StopID=%d"

	request, _ := http.NewRequest("GET", fmt.Sprintf(url, stopNumber), nil)
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("api_key", config.Api.Key)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("failed to hit API: %v", err)
	}
	defer response.Body.Close()

	// Read and print the response body
	body, _ := io.ReadAll(response.Body)

	var nbs NextBusServiceResponse
	err = json.Unmarshal(body, &nbs)
	if err != nil {
		fmt.Println("%v",err)
		return nil 
	}

	transitPredictions := make([]transit.Prediction, len(nbs.Predictions))

	for i, prediction := range nbs.Predictions {
		transitPredictions[i] = prediction
	} 

	//filterFunc := func()

	return lo.Filter(transitPredictions, func(p transit.Prediction, _ int) bool {
		return p.GetRouteID() == *busLineFilter
	})
}

type NextBusServiceResponse struct {
	StopName string `json:"StopName"`
	Predictions []*BusPrediction `json:"Predictions"`
}

type BusPrediction struct {
	RouteID string `json:"RouteID"`
	DirectionText string `json:"DirectionText"`
	Minutes int `json:"Minutes"`
}

func (b *BusPrediction) GetMinutes() string {
	return string(b.Minutes);
}

func (b *BusPrediction) GetDirectionText() string {
	return b.DirectionText;
}

func (b *BusPrediction) GetRouteID() string {
	return b.RouteID;
}

func (b BusPrediction) GetPrediction() string {
  return fmt.Sprintf("%s %s: %d\n", b.RouteID, b.DirectionText, b.Minutes)
}

