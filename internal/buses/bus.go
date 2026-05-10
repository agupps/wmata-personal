package buses

import (
	"net/http"
	"fmt"
	"log"
	"io"
	"encoding/json"
	"github.com/spf13/viper"
	"wmata/internal/transit"
)


func BusStops(stopNumber int) string {
	url := "https://api.wmata.com/NextBusService.svc/json/jPredictions?StopID=%d"

	apiKey := viper.GetString("api.key")

	request, _ := http.NewRequest("GET", fmt.Sprintf(url, stopNumber), nil)
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("api_key", apiKey)

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
		return ""
	}

	d72Predictions := []transit.Prediction{}

	for _, prediction := range nbs.Predictions {
		if prediction.RouteID == "D72" {
			d72Predictions = append(d72Predictions, prediction)
		}
	}

	return transit.PredictionsToString(d72Predictions)
}

//func predictionsToString(predictions []*BusPrediction) string {
	//out := ""
	//for _, p := range predictions {
	//	out += p.GetPrediction()
	//}
	//return out
//}


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

